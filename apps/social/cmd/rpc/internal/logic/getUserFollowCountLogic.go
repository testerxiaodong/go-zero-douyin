package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFollowCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFollowCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowCountLogic {
	return &GetUserFollowCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFollowCountLogic) GetUserFollowCount(in *pb.GetUserFollowCountReq) (*pb.GetUserFollowCountResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get user follow count with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get user follow count with empty user_id")
	}

	// 查询redis
	existsResult, err := l.svcCtx.Redis.ExistsCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, in.GetUserId()))
	if err != nil {
		logx.Errorf("get user follow count from redis key exist failed, err: %v, follower_id: %d", err, in.GetUserId())
	}

	// 在redis中有数据，直接返回
	if existsResult == true {
		val, err := l.svcCtx.Redis.ScardCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, in.GetUserId()))
		if err != nil {
			logx.Errorf("get user follow count from redis failed, err: %v, follower_id: %d", err, in.GetUserId())
		}
		// 更新缓存失效时间
		err = l.svcCtx.Redis.ExpireCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, in.GetUserId()), xconst.RedisExpireTime)
		if err != nil {
			logx.Errorf("set redis user follow count key expire time failed, err: %v, follower_id: %d", err, in.GetUserId())
		}
		return &pb.GetUserFollowCountResp{FollowCount: val}, nil
	}

	// 从数据库中获取数据
	key := cast.ToString(in.GetUserId())
	count, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.GetUserFollowCountFromDb(in.GetUserId())
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "get user follow count from mysql failed, err: %v, follower_id: %d", err, in.GetUserId())
	}

	// 类型断言，获得int64类型的数据
	countInt64, ok := count.(int64)
	if !ok {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "type assert failed")
	}

	// 异步构建缓存
	go l.BuildUserFollowCountCache(in.GetUserId())
	return &pb.GetUserFollowCountResp{FollowCount: countInt64}, nil
}

func (l *GetUserFollowCountLogic) GetUserFollowCountFromDb(userId int64) (int64, error) {
	followQuery := l.svcCtx.Query.Follow
	count, err := followQuery.WithContext(l.ctx).Where(followQuery.FollowerID.Eq(userId)).Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (l *GetUserFollowCountLogic) BuildUserFollowCountCache(followerId int64) {
	// 获取分布式锁的键
	lockKey := utils.GetRedisLockKeyWithPrefix(xconst.RedisBuildUserFollowCountCacheLockPrefix, followerId)
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)
	lock.SetExpire(1)

	// 获取分布式锁
	acquire, err := lock.Acquire()
	if err != nil {
		return
	}

	// 延迟释放分布式锁
	defer func(lock *redis.RedisLock) {
		_, err := lock.Release()
		if err != nil {

		}
	}(lock)

	// 获取成功，设置缓存以及失效时间
	if acquire {
		// 复制ctx，防止异步调用时logic.ctx结束
		ctx := contextx.ValueOnlyFrom(l.ctx)

		// 查询用户关注列表
		followQuery := l.svcCtx.Query.Follow
		follows, err := followQuery.WithContext(ctx).Where(followQuery.FollowerID.Eq(followerId)).Find()
		if err != nil {
			logx.WithContext(l.ctx).Errorf("find user follow list from db failed, follower_id: %d", followerId)
			return
		}
		if len(follows) > 0 {
			idList := make([]interface{}, 0, len(follows))
			for _, follow := range follows {
				idList = append(idList, follow.UserID)
			}
			_, err := l.svcCtx.Redis.SaddCtx(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, followerId), idList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis video liked by user cache  failed, err: %v", err)
				return
			}
			err = l.svcCtx.Redis.ExpireCtx(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, followerId), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set video liked by user redis key expire time failed, err: %v", err)
				return
			}
		}
	}
}

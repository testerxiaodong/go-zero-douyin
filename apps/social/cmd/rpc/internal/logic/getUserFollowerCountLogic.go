package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFollowerCountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFollowerCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowerCountLogic {
	return &GetUserFollowerCountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFollowerCountLogic) GetUserFollowerCount(in *pb.GetUserFollowerCountReq) (*pb.GetUserFollowerCountResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower count with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower count with empty user_id")
	}

	// 查询redis
	existsResult, err := l.svcCtx.Redis.Exists(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, in.GetUserId()))
	if err != nil {
		logx.Errorf("get user follower count from redis key exist failed, err: %v, user_id: %d", err, in.GetUserId())
	}

	// redis中有数据，直接返回
	if existsResult == true {
		val, err := l.svcCtx.Redis.Scard(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, in.GetUserId()))
		if err != nil {
			logx.Errorf("get user follower count from redis failed, err: %v, user_id: %d", err, in.GetUserId())
		}
		// 更新缓存失效时间
		err = l.svcCtx.Redis.Expire(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, in.GetUserId()), xconst.RedisExpireTime)
		if err != nil {
			logx.Errorf("set user follower count key expire time failed, err: %v, user_id: %d", err, in.GetUserId())
		}
		return &pb.GetUserFollowerCountResp{FollowerCount: val}, nil
	}

	// 从数据库中获取数据
	key := cast.ToString(in.GetUserId())
	count, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.svcCtx.FollowDo.GetUserFollowerCount(l.ctx, in.GetUserId())
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "get user follower count from mysql failed, err: %v, user_id: %v", err, in.GetUserId())
	}
	countInt64, ok := count.(int64)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("convert user follower count to int64 failed"), "user_id: %d", in.GetUserId())
	}

	// 异步构建缓存
	go l.BuildUserFollowerCountCache(in.GetUserId())

	return &pb.GetUserFollowerCountResp{FollowerCount: countInt64}, nil
}

func (l *GetUserFollowerCountLogic) BuildUserFollowerCountCache(userId int64) {
	// 获取分布式锁的键
	lockKey := utils.GetRedisLockKeyWithPrefix(xconst.RedisBuildUserFollowerCountCacheLockPrefix, userId)
	lock := l.svcCtx.Redis.NewRedisLock(lockKey)
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

	// 设置缓存以及失效时间
	if acquire {
		// 复制ctx，防止异步调用时logic.ctx结束
		ctx := contextx.ValueOnlyFrom(l.ctx)

		// 查询用户的粉丝列表
		ids, err := l.svcCtx.FollowDo.GetUserFollowerIdList(ctx, userId)
		if err != nil {
			logx.WithContext(ctx).Errorf("find user follower list from mysql failed, err: %v user_id: 5d", err, userId)
		}
		if len(ids) > 0 {
			idList := make([]interface{}, 0, len(ids))
			for _, id := range ids {
				idList = append(idList, id)
			}
			// set类型数据
			_, err := l.svcCtx.Redis.Sadd(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, userId), idList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis video liked by user cache  failed, err: %v", err)
				return
			}
			// 设置失效时间
			err = l.svcCtx.Redis.Expire(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, userId), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set video liked by user redis key expire time failed, err: %v", err)
				return
			}
		}
	}
}

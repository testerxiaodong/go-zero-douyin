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

type GetUserFollowIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFollowIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowIdListLogic {
	return &GetUserFollowIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFollowIdListLogic) GetUserFollowIdList(in *pb.GetUserFollowIdListReq) (*pb.GetUserFollowIdListResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follow id list with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follow id list with empty user_id")
	}

	// 查询redis
	existsResult, err := l.svcCtx.Redis.ExistsCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, in.GetUserId()))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("get redis user follow id list key exist failed"), "err: %v, follower_id: %d", err, in.GetUserId())
	}

	// redis中有数据，直接返回
	if existsResult == true {
		idList, err := l.svcCtx.Redis.SmembersCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, in.GetUserId()))
		if err != nil {
			// 查询失败，记录日志
			logx.WithContext(l.ctx).Errorf("get redis user follow id list failed, err: %v, follower_id: %d", err, in.GetUserId())
		}
		if len(idList) > 0 {
			resp := &pb.GetUserFollowIdListResp{UserIdList: make([]int64, 0)}
			for _, idStr := range idList {
				idInt64 := cast.ToInt64(idStr)
				resp.UserIdList = append(resp.UserIdList, idInt64)
			}
			// 更新缓存失效时间
			err := l.svcCtx.Redis.ExpireCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, in.GetUserId()), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("set redis user follow id list key expire time failed, err: %v, follower_id: %d", err, in.GetUserId())
			}
			// 返回数据
			return resp, nil
		}
	}

	// 从mysql中获取数据
	key := cast.ToString(in.GetUserId())
	idList, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.GetUserFollowIdListFromDb(in.GetUserId())
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "get user follow  id liet from mysql failed, err: %v", err)
	}

	// 类型断言
	idInt64List, ok := idList.([]int64)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("type assert failed, int64Slice"), "data: %v", idList)
	}

	// 异步构建缓存
	go l.BuildUserFollowIdListCache(in.GetUserId())

	return &pb.GetUserFollowIdListResp{UserIdList: idInt64List}, nil
}

func (l *GetUserFollowIdListLogic) GetUserFollowIdListFromDb(followerId int64) ([]int64, error) {
	followQuery := l.svcCtx.Query.Follow
	follows, err := followQuery.WithContext(l.ctx).Where(followQuery.FollowerID.Eq(followerId)).Find()
	if err != nil {
		return nil, err
	}
	if len(follows) > 0 {
		idList := make([]int64, 0)
		for _, follow := range follows {
			idList = append(idList, follow.UserID)
		}
		return idList, nil
	}
	return []int64{}, nil
}

func (l *GetUserFollowIdListLogic) BuildUserFollowIdListCache(followId int64) {
	// 获取分布式锁键
	lockKey := utils.GetRedisLockKeyWithPrefix(xconst.RedisBuildUserFollowCountCacheLockPrefix, followId)
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

	// 更新缓存
	if acquire {
		// 复制ctx，防止异步调用时logic.ctx结束
		ctx := contextx.ValueOnlyFrom(l.ctx)

		// 从数据库中查询视频点赞用户列表
		followQuery := l.svcCtx.Query.Follow
		follows, err := followQuery.WithContext(ctx).Where(followQuery.FollowerID.Eq(followId)).Find()
		if err != nil {
			logx.WithContext(ctx).Errorf("find user follow id list failed, err: %v", err)
		}
		if len(follows) > 0 {
			idList := make([]interface{}, 0, len(follows))
			for _, follow := range follows {
				idList = append(idList, follow.UserID)
			}

			// 设置缓存
			_, err := l.svcCtx.Redis.SaddCtx(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, followId), idList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis user follow id list cache  failed, err: %v", err)
				return
			}

			// 设置缓存失效时间
			err = l.svcCtx.Redis.ExpireCtx(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowUserPrefix, followId), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set redis user follow id list cache key expire time failed, err: %v", err)
				return
			}
		}
		return
	}
}

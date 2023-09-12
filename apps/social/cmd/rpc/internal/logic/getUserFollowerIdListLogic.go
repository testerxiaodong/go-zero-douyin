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

type GetUserFollowerIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserFollowerIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowerIdListLogic {
	return &GetUserFollowerIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserFollowerIdListLogic) GetUserFollowerIdList(in *pb.GetUserFollowerIdListReq) (*pb.GetUserFollowerIdListResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower id list with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user follower id list with empty user_id")
	}

	// 查询redis
	existsResult, err := l.svcCtx.Redis.Exists(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, in.GetUserId()))
	if err != nil {
		logx.WithContext(l.ctx).Errorf("get redis user follower id list key exist failed, err: %v, user_id: %d", err, in.GetUserId())
	}

	// redis中有数据，直接返回
	if existsResult == true {
		idList, err := l.svcCtx.Redis.Smembers(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, in.GetUserId()))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("get redis user follower id list failed, err: %v, user_id: %d", err, in.GetUserId())
		}
		if len(idList) > 0 {
			resp := &pb.GetUserFollowerIdListResp{UserIdList: make([]int64, 0)}
			for _, idStr := range idList {
				idInt64 := cast.ToInt64(idStr)
				resp.UserIdList = append(resp.UserIdList, idInt64)
			}
			// 更新缓存失效时间
			err := l.svcCtx.Redis.Expire(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, in.GetUserId()), xconst.RedisExpireTime)
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
		return l.svcCtx.FollowDo.GetUserFollowerIdList(l.ctx, in.GetUserId())
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "get user follower id list from mysql failed, err: %v", err)
	}

	// 类型断言
	idInt64List, ok := idList.([]int64)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("type assert failed, int64Slice"), "data: %v", idList)
	}

	// 异步构建缓存
	if len(idInt64List) > 0 {
		buildCache := func() {
			l.BuildUserFollowerIdListCache(in.GetUserId())
		}
		tsg := utils.NewTestGo()
		tsg.RunSafe(buildCache)
	}

	return &pb.GetUserFollowerIdListResp{UserIdList: idInt64List}, nil

}

func (l *GetUserFollowerIdListLogic) BuildUserFollowerIdListCache(userId int64) {
	// 获取分布式锁键
	lockKey := utils.GetRedisLockKeyWithPrefix(xconst.RedisBuildUserFollowerCountCacheLockPrefix, userId)
	lock := l.svcCtx.Redis.NewRedisLock(lockKey)
	lock.SetExpire(1)

	// 复制ctx，防止异步调用时logic.ctx结束
	ctx := contextx.ValueOnlyFrom(l.ctx)

	// 获取分布式锁
	acquire, err := lock.Acquire()
	if err != nil {
		logx.WithContext(ctx).Errorf("获取分布式锁失败，lockKey: %s, err: %v", lockKey, err)
		return
	}

	// 延迟释放分布式锁
	defer func(lock *redis.RedisLock) {
		_, err := lock.Release()
		if err != nil {
			logx.WithContext(ctx).Errorf("释放分布式锁失败，lockKey: %s, err: %v", lockKey, err)
		}
	}(lock)

	// 更新缓存
	if acquire {
		// 从数据库中查询视频点赞用户列表
		ids, err := l.svcCtx.FollowDo.GetUserFollowerIdList(ctx, userId)
		if err != nil {
			logx.WithContext(ctx).Errorf("find user follower id list failed, err: %v", err)
		}
		if len(ids) > 0 {
			idList := make([]interface{}, 0, len(ids))
			for _, id := range ids {
				idList = append(idList, id)
			}

			// 设置缓存
			_, err := l.svcCtx.Redis.Sadd(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, userId), idList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis user follower id list cache  failed, err: %v", err)
				return
			}

			// 设置缓存失效时间
			err = l.svcCtx.Redis.Expire(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserFollowedByUserPrefix, userId), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set redis user follower id list cache key expire time failed, err: %v", err)
				return
			}
		}
		return
	}
}

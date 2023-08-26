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

type GetUserLikeVideoIdListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLikeVideoIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLikeVideoIdListLogic {
	return &GetUserLikeVideoIdListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLikeVideoIdListLogic) GetUserLikeVideoIdList(in *pb.GetUserLikeVideoIdListReq) (*pb.GetUserLikeVideoIdListResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user like video id list with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user like video id list with empty user_id")
	}

	// 查询redis
	result, err := l.svcCtx.Redis.ExistsCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, in.GetUserId()))
	if err != nil {
		// 查询失败，记录日志
		logx.WithContext(l.ctx).Errorf("get redis user like video id list key exist failed, err: %v, user_id: %d", err, in.GetUserId())
	}
	// redis中key存在，直接返回redis数据
	if result == true {
		idSet, err := l.svcCtx.Redis.SmembersCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, in.GetUserId()))
		if err != nil {
			// 查询失败，记录日志
			logx.WithContext(l.ctx).Errorf("get redis user like video id list failed, err: %v, user_id: %d", err, in.GetUserId())
		}
		if len(idSet) > 0 {
			resp := &pb.GetUserLikeVideoIdListResp{VideoIdList: make([]int64, 0)}
			for _, idStr := range idSet {
				idInt64 := cast.ToInt64(idStr)
				resp.VideoIdList = append(resp.VideoIdList, idInt64)
			}
			// 更新缓存失效时间
			err := l.svcCtx.Redis.ExpireCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, in.GetUserId()), xconst.RedisExpireTime)
			if err != nil {
				// 设置缓存失效时间失败，记录日志
				logx.WithContext(l.ctx).Errorf("Set redis user like video id list key expire time failed, err: %v, user_id: %d", err, in.GetUserId())
			}
			return resp, nil
		}
	}

	// 查询mysql
	key := cast.ToString(in.GetUserId())
	idList, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.GetUserLikeVideoIdListFromDb(in.GetUserId())
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "get user like video id liet from mysql failed, err: %v", err)
	}

	// 类型断言
	idInt64List, ok := idList.([]int64)
	if !ok {
		return nil, errors.Wrapf(xerr.NewErrMsg("type assert failed, int64Slice"), "data: %v", idList)
	}

	// 异步构建缓存
	go l.BuildUserLikeVideoCache(in.GetUserId())

	return &pb.GetUserLikeVideoIdListResp{VideoIdList: idInt64List}, nil
}

func (l *GetUserLikeVideoIdListLogic) GetUserLikeVideoIdListFromDb(userId int64) ([]int64, error) {
	likeQuery := l.svcCtx.Query.Like
	likes, err := likeQuery.WithContext(l.ctx).Where(likeQuery.UserID.Eq(userId)).Find()
	if err != nil {
		return nil, err
	}
	if len(likes) > 0 {
		idList := make([]int64, 0)
		for _, like := range likes {
			idList = append(idList, like.VideoID)
		}
		return idList, nil
	}
	return []int64{}, nil
}

func (l *GetUserLikeVideoIdListLogic) BuildUserLikeVideoCache(userId int64) {
	// 获取分布式锁键
	lockKey := utils.GetRedisLockKeyWithPrefix(xconst.RedisBuildUserLikeVideoCacheLockPrefix, userId)
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
		likeQuery := l.svcCtx.Query.Like
		list, err := likeQuery.WithContext(ctx).Where(likeQuery.UserID.Eq(userId)).Find()
		if err != nil {
			logx.WithContext(ctx).Errorf("find user like video id list failed, err: %v", err)
		}
		if len(list) > 0 {
			idList := make([]interface{}, 0, len(list))
			for _, video := range list {
				idList = append(idList, video.VideoID)
			}

			// 设置缓存
			_, err := l.svcCtx.Redis.SaddCtx(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, userId), idList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis video liked by user cache  failed, err: %v", err)
				return
			}

			// 设置缓存失效时间
			err = l.svcCtx.Redis.ExpireCtx(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, userId), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set video liked by user redis key expire time failed, err: %v", err)
				return
			}
		}
		return
	}
}

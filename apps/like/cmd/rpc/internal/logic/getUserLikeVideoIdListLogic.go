package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-douyin/apps/like/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/like/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

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
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user like video is list with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get user like video is list with empty user_id")
	}

	// 查询redis
	result, err := l.svcCtx.Redis.ExistsCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, in.GetUserId()))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("get redis user like video is list key exist failed"), "user_id: %d, err: %v", in.GetUserId(), err)
	}
	// redis中key存在，直接返回redis数据
	if result == true {
		idSet, err := l.svcCtx.Redis.SmembersCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, in.GetUserId()))
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrMsg("get redis user like video is list failed"), "user_id: %d, err: %v", in.GetUserId(), err)
		}
		if len(idSet) > 0 {
			resp := &pb.GetUserLikeVideoIdListResp{VideoIdList: make([]int64, 0)}
			for _, idStr := range idSet {
				idInt64 := cast.ToInt64(idStr)
				resp.VideoIdList = append(resp.VideoIdList, idInt64)
			}
			return resp, nil
		}
		return nil, nil
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
	// 重新构建缓存
	l.BuildUserLikeVideoCache(in.GetUserId())
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
	return nil, nil
}

func (l *GetUserLikeVideoIdListLogic) BuildUserLikeVideoCache(userId int64) {
	lockKey := "build_user_like_video_id_list_key"
	lock := redis.NewRedisLock(l.svcCtx.Redis, lockKey)
	lock.SetExpire(1)
	acquire, err := lock.Acquire()
	if err != nil {
		return
	}
	defer func(lock *redis.RedisLock) {
		_, err := lock.Release()
		if err != nil {

		}
	}(lock)
	if acquire {
		likeQuery := l.svcCtx.Query.Like
		list, err := likeQuery.WithContext(l.ctx).Where(likeQuery.UserID.Eq(userId)).Find()
		if err != nil {
			logx.WithContext(l.ctx).Errorf("find user like video id list failed, err: %v", err)
		}
		if len(list) > 0 {
			idList := make([]interface{}, 0, len(list))
			for _, video := range list {
				idList = append(idList, video.VideoID)
			}
			_, err := l.svcCtx.Redis.SaddCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, userId), idList...)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("add redis video liked by user cache  failed, err: %v", err)
				return
			}
			err = l.svcCtx.Redis.ExpireCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisUserLikeVideoPrefix, userId), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(l.ctx).Errorf("set video liked by user redis key expire time failed, err: %v", err)
				return
			}
		}
		return
	}
}

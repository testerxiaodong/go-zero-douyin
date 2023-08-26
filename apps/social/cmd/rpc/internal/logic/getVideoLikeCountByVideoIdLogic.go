package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/contextx"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"
	"strconv"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoLikeCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoLikeCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoLikeCountByVideoIdLogic {
	return &GetVideoLikeCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoLikeCountByVideoIdLogic) GetVideoLikeCountByVideoId(in *pb.GetVideoLikeCountByVideoIdReq) (*pb.GetVideoLikeCountByVideoIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video like count with empty param")
	}
	if in.GetVideoId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video like count with empty video_id")
	}

	// 先查询redis
	result, err := l.svcCtx.Redis.ExistsCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, in.GetVideoId()))
	if err != nil {
		// 查询失败。记录日志
		logc.Errorf(l.ctx, "get redis video like count failed, err: %v video_id: %d", err, in.GetVideoId())
	}

	// redis有数据，直接返回
	if result == true {
		val, err := l.svcCtx.Redis.ScardCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, in.GetVideoId()))
		if err != nil {
			// 查询失败。记录日志
			logc.Errorf(l.ctx, "get redis video like count failed, err: %v video_id: %d", err, in.GetVideoId())
		}
		return &pb.GetVideoLikeCountByVideoIdResp{LikeCount: val}, nil
	}

	// 查询数据库
	key := strconv.Itoa(int(in.GetVideoId()))
	count, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.GetVideoLikedCountFromDb(in.GetVideoId())
	})
	if err != nil {
		return nil, err
	}
	countInt64, ok := count.(int64)
	if !ok {
		return &pb.GetVideoLikeCountByVideoIdResp{}, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "type assert failed")
	}

	// 异步构建缓存
	go l.BuildVideoLikedByUserCache(in.GetVideoId())

	return &pb.GetVideoLikeCountByVideoIdResp{LikeCount: countInt64}, nil
}

// GetVideoLikedCountFromDb 从mysql中获取视频点赞数
func (l *GetVideoLikeCountByVideoIdLogic) GetVideoLikedCountFromDb(videoId int64) (int64, error) {
	likeQuery := l.svcCtx.Query.Like
	count, err := likeQuery.WithContext(l.ctx).Where(likeQuery.VideoID.Eq(videoId)).Count()
	if err != nil {
		return 0, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "get video like count from db failed, err: %v", err)
	}
	return count, nil
}

// BuildVideoLikedByUserCache 构建视频被点赞缓存
func (l *GetVideoLikeCountByVideoIdLogic) BuildVideoLikedByUserCache(videoId int64) {
	// 获取分布式的键
	lockKey := utils.GetRedisLockKeyWithPrefix(xconst.RedisBuildVideoLikedByUserCacheLockPrefix, videoId)
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
		// 复制logic.ctx，防止异步调用时父context结束
		ctx := contextx.ValueOnlyFrom(l.ctx)

		// 查询点赞视频的用户列表
		likeQuery := l.svcCtx.Query.Like
		list, err := likeQuery.WithContext(ctx).Where(likeQuery.VideoID.Eq(videoId)).Find()
		if err != nil {
			logx.WithContext(ctx).Errorf("find video liked by user with video_id failed, err: %v", err)
		}
		if len(list) > 0 {
			idList := make([]interface{}, 0, len(list))
			for _, video := range list {
				idList = append(idList, video.UserID)
			}

			// 设置缓存
			_, err := l.svcCtx.Redis.SaddCtx(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, videoId), idList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis video liked by user cache  failed, err: %v", err)
				return
			}

			// 设置缓存失效时间
			err = l.svcCtx.Redis.ExpireCtx(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, videoId), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set video liked by user redis key expire time failed, err: %v", err)
				return
			}
		}
		return
	}
}

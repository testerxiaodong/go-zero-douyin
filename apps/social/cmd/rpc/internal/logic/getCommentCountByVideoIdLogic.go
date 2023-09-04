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
	"strconv"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountByVideoIdLogic {
	return &GetCommentCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentCountByVideoIdLogic) GetCommentCountByVideoId(in *pb.GetCommentCountByVideoIdReq) (*pb.GetCommentCountByVideoIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video comment count with empty param")
	}
	if in.GetVideoId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video comment count with empty video_id")
	}

	// 从redis中获取数据
	result, err := l.svcCtx.Redis.Exists(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, in.GetVideoId()))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "get redis video comment count key exist failed: %v", err)
	}

	// redis中有数据，直接返回
	if result == true {
		val, err := l.svcCtx.Redis.Get(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, in.GetVideoId()))
		if err != nil {
			logx.WithContext(l.ctx).Errorf("get redis video comment count failed: %v, video_id: %d", err, in.GetVideoId())
		}
		count := cast.ToInt64(val)
		err = l.svcCtx.Redis.Expire(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, in.GetVideoId()), xconst.RedisExpireTime)
		if err != nil {
			logx.WithContext(l.ctx).Errorf("set redis video comment count key expire time failed: %v, video_id: %d", err, in.GetVideoId())
		}
		return &pb.GetCommentCountByVideoIdResp{Count: count}, nil
	}

	// 从mysql中获取数据
	key := strconv.Itoa(int(in.GetVideoId()))
	count, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.GetVideoCommentFromDb(in.GetVideoId())
	})
	if err != nil {
		return nil, err
	}
	countInt64, ok := count.(int64)
	if !ok {
		return &pb.GetCommentCountByVideoIdResp{}, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "type assert failed")
	}

	// 异步构建缓存
	go l.BuildVideoCommentCountCache(in.VideoId, countInt64)

	return &pb.GetCommentCountByVideoIdResp{Count: countInt64}, nil
}

func (l *GetCommentCountByVideoIdLogic) GetVideoCommentFromDb(videoId int64) (int64, error) {
	count, err := l.svcCtx.CommentDo.GetCommentCountByVideoId(l.ctx, videoId)
	if err != nil {
		return 0, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "get mysql video commnet count failed: %v", err)
	}
	return count, nil
}

func (l *GetCommentCountByVideoIdLogic) BuildVideoCommentCountCache(videoId int64, commentCount int64) {
	// 获取分布式锁的键
	lockKey := utils.GetRedisLockKeyWithPrefix(xconst.RedisBuildVideoCommentCountCacheLockPrefix, videoId)
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

	// 获取成功，设置缓存以及失效时间
	if acquire {
		// 复制logic.ctx，防止异步更新缓存时，父context结束了
		ctx := contextx.ValueOnlyFrom(l.ctx)

		// 设置缓存
		commentCountStr := strconv.Itoa(int(commentCount))
		err = l.svcCtx.Redis.Set(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, videoId), commentCountStr)
		if err != nil {
			logx.WithContext(ctx).Errorf("set video comment cache failed, video_id: %d", videoId)
			return
		}

		// 设置缓存失效时间
		err := l.svcCtx.Redis.Expire(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, videoId), xconst.RedisExpireTime)
		if err != nil {
			logx.WithContext(ctx).Errorf("set video comment cache expire time failed, video_id: %d", videoId)
			return
		}
	}
}

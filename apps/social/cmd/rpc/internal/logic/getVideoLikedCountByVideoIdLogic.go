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

type GetVideoLikedCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoLikedCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoLikedCountByVideoIdLogic {
	return &GetVideoLikedCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoLikedCountByVideoIdLogic) GetVideoLikedCountByVideoId(in *pb.GetVideoLikedCountByVideoIdReq) (*pb.GetVideoLikedCountByVideoIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video like count with empty param")
	}
	if in.GetVideoId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video like count with empty video_id")
	}

	// 先查询redis
	result, err := l.svcCtx.Redis.Exists(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, in.GetVideoId()))
	if err != nil {
		// 查询失败。记录日志
		logc.Errorf(l.ctx, "get redis video like count failed, err: %v video_id: %d", err, in.GetVideoId())
	}

	// redis有数据，直接返回
	if result == true {
		val, err := l.svcCtx.Redis.Scard(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, in.GetVideoId()))
		if err != nil {
			// 查询失败。记录日志
			logc.Errorf(l.ctx, "get redis video like count failed, err: %v video_id: %d", err, in.GetVideoId())
		}
		// 更新缓存失效时间
		err = l.svcCtx.Redis.Expire(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, in.GetVideoId()), xconst.RedisExpireTime)
		if err != nil {
			// 设置缓存失效时间失败，记录日志
			logx.WithContext(l.ctx).Errorf("Set redis user like video id list key expire time failed, err: %v, user_id: %d", err, in.GetVideoId())
		}
		return &pb.GetVideoLikedCountByVideoIdResp{LikeCount: val}, nil
	}

	// 查询数据库
	key := strconv.Itoa(int(in.GetVideoId()))
	count, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.svcCtx.LikeDo.GetVideoLikedCount(l.ctx, in.GetVideoId())
	})
	if err != nil {
		return nil, err
	}
	countInt64, ok := count.(int64)
	if !ok {
		return &pb.GetVideoLikedCountByVideoIdResp{}, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "type assert failed")
	}

	// 异步构建缓存
	go l.BuildVideoLikedByUserCache(in.GetVideoId())

	return &pb.GetVideoLikedCountByVideoIdResp{LikeCount: countInt64}, nil
}

// BuildVideoLikedByUserCache 构建视频被点赞缓存
func (l *GetVideoLikedCountByVideoIdLogic) BuildVideoLikedByUserCache(videoId int64) {
	// 获取分布式的键
	lockKey := utils.GetRedisLockKeyWithPrefix(xconst.RedisBuildVideoLikedByUserCacheLockPrefix, videoId)
	lock := l.svcCtx.Redis.NewRedisLock(lockKey)
	lock.SetExpire(1)

	// 复制logic.ctx，防止异步调用时父context结束
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
		// 查询点赞视频的用户列表
		ids, err := l.svcCtx.LikeDo.GetVideoLikedByUserIdList(l.ctx, videoId)
		if err != nil {
			logx.WithContext(ctx).Errorf("find video liked by user with video_id failed, err: %v", err)
		}
		if len(ids) > 0 {
			idList := make([]interface{}, 0, len(ids))
			for _, id := range ids {
				idList = append(idList, id)
			}

			// 设置缓存
			_, err := l.svcCtx.Redis.Sadd(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, videoId), idList...)
			if err != nil {
				logx.WithContext(ctx).Errorf("add redis video liked by user cache  failed, err: %v", err)
				return
			}

			// 设置缓存失效时间
			err = l.svcCtx.Redis.Expire(ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, videoId), xconst.RedisExpireTime)
			if err != nil {
				logx.WithContext(ctx).Errorf("set video liked by user redis key expire time failed, err: %v", err)
				return
			}
		}
		return
	}
}

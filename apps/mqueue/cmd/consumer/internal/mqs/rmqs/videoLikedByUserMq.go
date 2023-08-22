package rmqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	messageType "go-zero-douyin/common/message"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
)

// VideoLikedByUserMq 视频评论缓存信息
type VideoLikedByUserMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoLikedByUserMq(ctx context.Context, svcCtx *svc.ServiceContext) *VideoLikedByUserMq {
	return &VideoLikedByUserMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *VideoLikedByUserMq) Consume(message string) error {
	// 获取消息内容：视频id
	var videoLikedByUserMessage messageType.VideoLikedByUserMessage
	if err := json.Unmarshal([]byte(message), &videoLikedByUserMessage); err != nil {
		logx.WithContext(v.ctx).Error("videoLikedByUserMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	// 删除视频缓存
	_, err := v.svcCtx.Redis.Del(utils.GetRedisKeyWithPrefix(xconst.RedisVideoLikedByUserPrefix, videoLikedByUserMessage.VideoId))
	if err != nil {
		// 少于重试最高次数，重新入队
		return err
	}
	return nil
}

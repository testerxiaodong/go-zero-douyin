package rmqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	messageType "go-zero-douyin/common/message"
	"go-zero-douyin/common/xconst"
)

// VideoSectionMq 视频评论缓存信息
type VideoSectionMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoSectionMq(ctx context.Context, svcCtx *svc.ServiceContext) *VideoSectionMq {
	return &VideoSectionMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *VideoSectionMq) Consume(message string) error {
	// 获取消息内容：视频id
	var videoSectionMessage messageType.VideoSectionMessage
	if err := json.Unmarshal([]byte(message), &videoSectionMessage); err != nil {
		logx.WithContext(v.ctx).Error("videoSectionMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	// 删除视频缓存
	_, err := v.svcCtx.Redis.Del(xconst.RedisVideoSection)
	if err != nil {
		// 少于重试最高次数，重新入队
		return err
	}
	return nil
}

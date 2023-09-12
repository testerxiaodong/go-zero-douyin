package rmqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	messageType "go-zero-douyin/common/message"
	"go-zero-douyin/common/xconst"
)

// VideoTagMq 视频评论缓存信息
type VideoTagMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoTagMq(ctx context.Context, svcCtx *svc.ServiceContext) *VideoTagMq {
	return &VideoTagMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *VideoTagMq) Consume(message string) error {
	// 获取消息内容：视频id
	var videoTagMessage messageType.VideoTagMessage
	if err := json.Unmarshal([]byte(message), &videoTagMessage); err != nil {
		logx.WithContext(v.ctx).Error("videoTagMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	// 删除视频缓存
	_, err := v.svcCtx.Redis.Del(xconst.RedisVideoTag)
	if err != nil {
		// 少于重试最高次数，重新入队
		return err
	}
	return nil
}

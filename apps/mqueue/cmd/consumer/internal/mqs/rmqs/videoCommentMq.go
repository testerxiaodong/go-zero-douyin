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

// VideoCommentMq 视频评论缓存信息
type VideoCommentMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoCommentMq(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCommentMq {
	return &VideoCommentMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *VideoCommentMq) Consume(message string) error {
	// 获取消息内容：视频id
	var videoCommentMessage messageType.VideoCommentMessage
	if err := json.Unmarshal([]byte(message), &videoCommentMessage); err != nil {
		logx.WithContext(v.ctx).Error("videoCommentMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	// 删除视频缓存
	_, err := v.svcCtx.Redis.Del(utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, videoCommentMessage.VideoId))
	if err != nil {
		// 少于重试最高次数，重新入队
		return err
	}
	return nil
}

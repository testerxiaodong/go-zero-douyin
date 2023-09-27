package kqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/message"
)

type VideoLikeMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoLikeMq(ctx context.Context, svcCtx *svc.ServiceContext) *VideoLikeMq {
	return &VideoLikeMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoLikeMq) Consume(_, val string) error {

	var videoLikeMessage message.VideoLikeMessage
	if err := json.Unmarshal([]byte(val), &videoLikeMessage); err != nil {
		logx.WithContext(l.ctx).Error("VideoLikeMq->Consume Unmarshal err : %v , val : %s", err, val)
		return err
	}
	if err := l.execService(videoLikeMessage); err != nil {
		logx.WithContext(l.ctx).Error("VideoLikeMq->execService  err : %v , val : %s , message:%+v", err, val, videoLikeMessage)
		return err
	}

	return nil
}

// 调用Rpc执行视频点赞业务
func (l *VideoLikeMq) execService(message message.VideoLikeMessage) error {
	_, err := l.svcCtx.SocialRpc.VideoLike(l.ctx, &pb.VideoLikeReq{UserId: message.UserId, VideoId: message.VideoId})
	if err != nil {
		return err
	}
	return nil
}

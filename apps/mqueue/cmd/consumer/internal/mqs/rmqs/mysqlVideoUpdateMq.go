package rmqs

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
	searchPb "go-zero-douyin/apps/search/cmd/rpc/pb"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	messageType "go-zero-douyin/common/message"
)

// MysqlVideoUpdateMq 用户关注缓存
type MysqlVideoUpdateMq struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMysqlVideoUpdateMq(ctx context.Context, svcCtx *svc.ServiceContext) *MysqlVideoUpdateMq {
	return &MysqlVideoUpdateMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (v *MysqlVideoUpdateMq) Consume(message string) error {
	// 获取消息内容：视频id
	var mysqlVideoUpdateMessage messageType.MysqlVideoUpdateMessage
	if err := json.Unmarshal([]byte(message), &mysqlVideoUpdateMessage); err != nil {
		logx.WithContext(v.ctx).Error("mysqlVideoUpdateMessage->Consume Unmarshal err : %v , val : %s", err, message)
		return err
	}
	videoInfoResp, err := v.svcCtx.VideoRpc.GetVideoById(v.ctx, &pb.GetVideoByIdReq{Id: mysqlVideoUpdateMessage.VideoId})
	if err != nil {
		return err
	}
	commentCountResp, err := v.svcCtx.SocialRpc.GetCommentCountByVideoId(v.ctx, &socialPb.GetCommentCountByVideoIdReq{VideoId: mysqlVideoUpdateMessage.VideoId})
	if err != nil {
		return err
	}
	likedCountResp, err := v.svcCtx.SocialRpc.GetVideoLikedCountByVideoId(v.ctx, &socialPb.GetVideoLikedCountByVideoIdReq{VideoId: mysqlVideoUpdateMessage.VideoId})
	if err != nil {
		return err
	}
	_, err = v.svcCtx.SearchRpc.SyncVideoInfo(v.ctx, &searchPb.SyncVideoInfoReq{Video: &searchPb.Video{
		Id:           videoInfoResp.Video.Id,
		Title:        videoInfoResp.Video.Title,
		SectionId:    videoInfoResp.Video.SectionId,
		Tags:         videoInfoResp.Video.Tags,
		OwnerId:      videoInfoResp.Video.OwnerId,
		OwnerName:    videoInfoResp.Video.OwnerName,
		PlayUrl:      videoInfoResp.Video.PlayUrl,
		CoverUrl:     videoInfoResp.Video.CoverUrl,
		CommentCount: commentCountResp.Count,
		LikeCount:    likedCountResp.LikeCount,
		CreateTime:   videoInfoResp.Video.CreateTime,
		UpdateTime:   videoInfoResp.Video.UpdateTime,
	}})
	if err != nil {
		return err
	}
	return nil
}

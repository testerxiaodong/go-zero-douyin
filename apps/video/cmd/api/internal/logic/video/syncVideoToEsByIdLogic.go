package video

import (
	"context"
	"github.com/pkg/errors"
	searchPb "go-zero-douyin/apps/search/cmd/rpc/pb"
	socialPb "go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncVideoToEsByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncVideoToEsByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncVideoToEsByIdLogic {
	return &SyncVideoToEsByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncVideoToEsByIdLogic) SyncVideoToEsById(req *types.SyncVideoToEsByIdReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}
	// 调用videoRpc获取视频信息
	videoInfo, err := l.svcCtx.VideoRpc.GetVideoById(l.ctx, &pb.GetVideoByIdReq{Id: req.VideoId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	// 调用socialRpc获取视频评论数
	videoCommentCountResp, err := l.svcCtx.SocialRpc.GetCommentCountByVideoId(l.ctx, &socialPb.GetCommentCountByVideoIdReq{VideoId: req.VideoId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	// 调用socialRpc获取视频点赞数
	videoLikeCountResp, err := l.svcCtx.SocialRpc.GetVideoLikedCountByVideoId(l.ctx, &socialPb.GetVideoLikedCountByVideoIdReq{VideoId: req.VideoId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	// 调用SearchRpc同步视频信息
	_, err = l.svcCtx.SearchRpc.SyncVideoInfo(l.ctx, &searchPb.SyncVideoInfoReq{Video: &searchPb.Video{
		Id:           videoInfo.Video.Id,
		Title:        videoInfo.Video.Title,
		SectionId:    videoInfo.Video.SectionId,
		Tags:         videoInfo.Video.Tags,
		OwnerId:      videoInfo.Video.OwnerId,
		PlayUrl:      videoInfo.Video.PlayUrl,
		CoverUrl:     videoInfo.Video.CoverUrl,
		CommentCount: videoCommentCountResp.Count,
		LikeCount:    videoLikeCountResp.LikeCount,
		CreateTime:   videoInfo.Video.CreateTime,
		UpdateTime:   videoInfo.Video.UpdateTime,
	}})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	return nil
}

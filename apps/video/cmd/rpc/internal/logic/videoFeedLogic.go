package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoFeedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVideoFeedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoFeedLogic {
	return &VideoFeedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VideoFeedLogic) VideoFeed(in *pb.VideoFeedReq) (*pb.VideoFeedResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get video feed with empty param")
	}
	if in.GetLastTimeStamp() == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get video feed with empty timestamp")
	}
	// 获取数据库数据
	lastTimeStamp := in.GetLastTimeStamp().GetSeconds()
	videoQuery := l.svcCtx.Query.Video
	videos, err := videoQuery.WithContext(l.ctx).Where(videoQuery.CreateTime.Lt(lastTimeStamp)).Limit(xconst.VideoFeedCount).Order(videoQuery.CreateTime.Desc()).Find()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Get video feed by last timestmap failed, err: %v", err)
	}
	// 拼接数据
	if len(videos) == 0 {
		return &pb.VideoFeedResp{}, nil
	}
	videosResp := &pb.VideoFeedResp{Videos: make([]*pb.Video, 0)}
	for _, video := range videos {
		singleVideoResp := &pb.Video{}
		singleVideoResp.Id = video.ID
		singleVideoResp.Title = video.Title
		singleVideoResp.OwnerId = video.OwnerID
		singleVideoResp.PlayUrl = video.PlayURL
		singleVideoResp.CoverUrl = video.CoverURL
		videosResp.Videos = append(videosResp.Videos, singleVideoResp)
	}
	return videosResp, nil
}

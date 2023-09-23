package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

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
	if in.GetLastTimeStamp() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get video feed with empty timestamp")
	}

	// 获取数据库数据
	lastTime := utils.FromUnixTimestampToTime(in.GetLastTimeStamp())
	builder := l.svcCtx.VideoModel.SelectBuilder().Where(squirrel.Lt{"create_time": lastTime}).Where(squirrel.Eq{"section_id": in.GetSectionId()})
	videos, err := l.svcCtx.VideoModel.FindPageListByPage(l.ctx, builder, 1, xconst.VideoFeedCount, "create_time DESC")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Get video feed by last timestmap failed, err: %v", err)
	}

	// 拼接数据
	if len(videos) == 0 {
		return &pb.VideoFeedResp{}, nil
	}
	videosResp := &pb.VideoFeedResp{Videos: make([]*pb.VideoInfo, 0)}
	for _, video := range videos {
		singleVideoResp := &pb.VideoInfo{}
		singleVideoResp.Id = video.Id
		singleVideoResp.Title = video.Title
		singleVideoResp.SectionId = video.SectionId
		singleVideoResp.TagIds = video.TagIds
		singleVideoResp.OwnerId = video.OwnerId
		singleVideoResp.OwnerName = video.OwnerName
		singleVideoResp.PlayUrl = video.PlayUrl
		singleVideoResp.CoverUrl = video.CoverUrl
		singleVideoResp.CreateTime = video.CreateTime.Unix()
		singleVideoResp.UpdateTime = video.UpdateTime.Unix()
		videosResp.Videos = append(videosResp.Videos, singleVideoResp)
	}
	return videosResp, nil
}

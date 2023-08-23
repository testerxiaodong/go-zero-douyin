package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserVideoListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserVideoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserVideoListLogic {
	return &UserVideoListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserVideoListLogic) UserVideoList(in *pb.UserVideoListReq) (*pb.UserVideoListResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get User Video List with empty param")
	}
	if in.GetUserId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Get user video list with empty user_id")
	}
	// 查询数据库数据
	videoQuery := l.svcCtx.Query.Video
	videos, err := videoQuery.WithContext(l.ctx).Where(videoQuery.OwnerID.Eq(in.GetUserId())).Find()
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Get user video list by user_id failed, err: %v", err)
	}
	// 没有数据，直接返回
	if len(videos) == 0 {
		return &pb.UserVideoListResp{}, nil
	}
	// 拼接数据
	resp := &pb.UserVideoListResp{Videos: make([]*pb.VideoInfo, 0)}
	for _, video := range videos {
		single := &pb.VideoInfo{}
		single.Id = video.ID
		single.OwnerId = video.OwnerID
		single.Title = video.Title
		single.PlayUrl = video.PlayURL
		single.CoverUrl = video.CoverURL
		resp.Videos = append(resp.Videos, single)
	}
	return resp, nil
}

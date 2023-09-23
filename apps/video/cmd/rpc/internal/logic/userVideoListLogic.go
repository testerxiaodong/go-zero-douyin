package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

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
	builder := l.svcCtx.VideoModel.SelectBuilder().Where(squirrel.Eq{"owner_id": in.GetUserId()})
	videos, total, err := l.svcCtx.VideoModel.FindPageListByPageWithTotal(l.ctx, builder, in.GetPage(), in.GetPageSize(), "create_time DESC")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Get user video list by user_id failed, err: %v", err)
	}

	// 没有数据，直接返回
	if len(videos) == 0 {
		return &pb.UserVideoListResp{}, nil
	}

	// 拼接数据
	resp := &pb.UserVideoListResp{Videos: make([]*pb.VideoInfo, 0)}
	resp.Total = total
	for _, video := range videos {
		single := &pb.VideoInfo{}
		single.Id = video.Id
		single.OwnerId = video.OwnerId
		single.OwnerName = video.OwnerName
		single.SectionId = video.SectionId
		single.TagIds = video.TagIds
		single.Title = video.Title
		single.PlayUrl = video.PlayUrl
		single.CoverUrl = video.CoverUrl
		single.CreateTime = video.CreateTime.Unix()
		single.UpdateTime = video.UpdateTime.Unix()
		resp.Videos = append(resp.Videos, single)
	}
	return resp, nil
}

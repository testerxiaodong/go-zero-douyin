package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrVideoNotFound = xerr.NewErrMsg("视频不存在")

type DeleteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVideoLogic {
	return &DeleteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteVideoLogic) DeleteVideo(in *pb.DeleteVideoReq) (*pb.DeleteVideoResp, error) {
	// todo: add your logic here and delete this line
	// 校验参数
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Delete video with empty param")
	}
	if in.GetUserId() == 0 || in.GetVideoId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Delete video with empty user_id or video_id")
	}

	// 查询视频是否是该用户发布的
	video, err := l.svcCtx.VideoModel.FindOne(l.ctx, in.GetVideoId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "Find video by id failed: %v", err)
	}

	// 视频不存在
	if video == nil {
		return nil, errors.Wrapf(ErrVideoNotFound, "video_id: %d", in.GetVideoId())
	}

	// 视频非该用户发布
	if video.OwnerId != in.GetUserId() {
		return nil, errors.Wrapf(xerr.NewErrMsg("视频非该用户发布，用户无权操作"), "video_id: %d", in.GetVideoId())
	}

	// 删除视频
	err = l.svcCtx.VideoModel.DeleteSoft(l.ctx, nil, video)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "delete video by id failed: %v", err)
	}

	// 返回响应
	return &pb.DeleteVideoResp{}, nil
}

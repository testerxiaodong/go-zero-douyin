package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/internal/model"
	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoLikedCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoLikedCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoLikedCountByVideoIdLogic {
	return &GetVideoLikedCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoLikedCountByVideoIdLogic) GetVideoLikedCountByVideoId(in *pb.GetVideoLikedCountByVideoIdReq) (*pb.GetVideoLikedCountByVideoIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video like count with empty param")
	}
	if in.GetVideoId() == 0 {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video like count with empty video_id")
	}
	// 查询数据库
	likeCount, err := l.svcCtx.LikeCountModel.FindOneByVideoId(l.ctx, in.GetVideoId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"查询视频点赞数失败, err: %v, video_id: %d", err, in.GetVideoId())
	}
	if likeCount == nil {
		return &pb.GetVideoLikedCountByVideoIdResp{LikeCount: 0}, nil
	}
	return &pb.GetVideoLikedCountByVideoIdResp{LikeCount: likeCount.LikeCount}, nil
}

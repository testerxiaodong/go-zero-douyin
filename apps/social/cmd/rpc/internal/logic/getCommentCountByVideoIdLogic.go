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

type GetCommentCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountByVideoIdLogic {
	return &GetCommentCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentCountByVideoIdLogic) GetCommentCountByVideoId(in *pb.GetCommentCountByVideoIdReq) (*pb.GetCommentCountByVideoIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video comment count with empty param")
	}
	if in.GetVideoId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video comment count with empty video_id")
	}
	// 查询数据
	commentCount, err := l.svcCtx.CommentCountModel.FindOneByVideoId(l.ctx, in.GetVideoId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询视频评论数失败, err: %v, video_id: %d", err, in.GetVideoId())
	}
	// 没有数据，说明评论数为零
	if commentCount == nil {
		return &pb.GetCommentCountByVideoIdResp{Count: 0}, nil
	}
	// 返回查询到的评论数
	return &pb.GetCommentCountByVideoIdResp{Count: commentCount.CommentCount}, nil
}

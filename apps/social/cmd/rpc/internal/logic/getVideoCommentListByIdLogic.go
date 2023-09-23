package logic

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetVideoCommentListByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetVideoCommentListByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVideoCommentListByIdLogic {
	return &GetVideoCommentListByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetVideoCommentListByIdLogic) GetVideoCommentListById(in *pb.GetCommentListByIdReq) (*pb.GetCommentListByIdResp, error) {
	// todo: add your logic here and delete this line
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get video comment list by id with empty param")
	}
	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "Get video comment list by id with empty id")
	}

	// 查询数据库
	builder := l.svcCtx.CommentModel.SelectBuilder().Where(squirrel.Eq{"video_id": in.GetId()})
	comments, n, err := l.svcCtx.CommentModel.FindPageListByPageWithTotal(l.ctx, builder, in.GetPage(), in.GetPageSize(), "create_time DESC")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "db get comment list by video_id failed: %v", err)
	}

	if len(comments) == 0 {
		return &pb.GetCommentListByIdResp{}, nil
	}

	resp := &pb.GetCommentListByIdResp{Comments: make([]*pb.Comment, 0)}
	resp.Total = n
	_ = copier.Copy(&resp.Comments, comments)
	return resp, nil
}

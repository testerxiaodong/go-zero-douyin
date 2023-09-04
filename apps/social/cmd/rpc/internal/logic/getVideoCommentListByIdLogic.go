package logic

import (
	"context"
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
	comments, err := l.svcCtx.CommentDo.GetCommentListByVideoId(l.ctx, in.GetId())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "db get comment list by video_id failed: %v", err)
	}

	if len(comments) == 0 {
		return &pb.GetCommentListByIdResp{}, nil
	}

	resp := &pb.GetCommentListByIdResp{Comments: make([]*pb.Comment, 0)}

	// 拼接数据
	for _, comment := range comments {
		single := &pb.Comment{}
		single.Id = comment.ID
		single.VideoId = comment.VideoID
		single.UserId = comment.UserID
		single.Content = comment.Content
		resp.Comments = append(resp.Comments, single)
	}

	return resp, nil
}

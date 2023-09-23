package comment

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VideoCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVideoCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VideoCommentListLogic {
	return &VideoCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VideoCommentListLogic) VideoCommentList(req *types.GetVideoCommentListReq) (resp *types.GetVideoCommentListResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 调用commentrpc
	comments, err := l.svcCtx.SocialRpc.GetVideoCommentListById(l.ctx,
		&pb.GetCommentListByIdReq{Id: req.VideoId, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	// 拼接响应
	resp = &types.GetVideoCommentListResp{Comments: make([]*types.Comment, 0)}
	err = copier.Copy(resp, comments)
	if err != nil {
		return nil, errors.Wrapf(err, "copy rpc resp to api failed, data: %v", comments)
	}

	return resp, nil
}

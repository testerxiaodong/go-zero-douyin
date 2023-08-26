package comment

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelCommentLogic {
	return &DelCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelCommentLogic) DelComment(req *types.DelCommentReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.ValidateZh(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}

	// 获取当前用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 调用commentrpc
	_, err := l.svcCtx.SocialRpc.DelComment(l.ctx, &pb.DelCommentReq{UserId: uid, CommentId: req.CommentId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}

	return nil
}

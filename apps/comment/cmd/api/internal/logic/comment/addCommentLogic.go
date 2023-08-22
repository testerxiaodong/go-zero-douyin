package comment

import (
	"context"
	commentPb "go-zero-douyin/apps/comment/cmd/rpc/pb"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/comment/cmd/api/internal/svc"
	"go-zero-douyin/apps/comment/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCommentLogic) AddComment(req *types.AddCommentReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := utils.GetValidator().ValidateZh(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}

	// 调用videorpc判断视频是否存在
	_, err := l.svcCtx.VideoRpc.GetVideoById(l.ctx, &pb.GetVideoByIdReq{Id: req.VideoId})
	if err != nil {
		return xerr.NewErrMsg("视频不存在")
	}

	// 获取当前用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 调用commentrpc
	_, err = l.svcCtx.CommentRpc.AddComment(l.ctx, &commentPb.AddCommentReq{VideoId: req.VideoId, UserId: uid, Content: req.Content})
	if err != nil {
		return err
	}
	return nil
}

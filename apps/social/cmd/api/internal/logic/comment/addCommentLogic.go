package comment

import (
	"context"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	pbVideo "go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"

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
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}

	// 调用 video rpc 判断视频是否存在
	_, err := l.svcCtx.VideoRpc.GetVideoById(l.ctx, &pbVideo.GetVideoByIdReq{Id: req.VideoId})
	if err != nil {
		return err
	}

	// 获取当前用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 调用 social rpc
	_, err = l.svcCtx.SocialRpc.AddComment(l.ctx, &pb.AddCommentReq{VideoId: req.VideoId, UserId: uid, Content: req.Content})
	if err != nil {
		return err
	}
	return nil
}

package like

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	pbVideo "go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelLikeLogic {
	return &DelLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelLikeLogic) DelLike(req *types.VideoUnlikeReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := utils.GetValidator().ValidateZh(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}

	// 调用videorpc查询视频是否存在
	_, err := l.svcCtx.VideoRpc.GetVideoById(l.ctx, &pbVideo.GetVideoByIdReq{Id: req.VideoId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}

	// 获取当前用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 调用likerpc
	_, err = l.svcCtx.SocialRpc.VideoUnlike(l.ctx, &pb.VideoUnlikeReq{VideoId: req.VideoId, UserId: uid})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}

	return nil
}

package like

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	pbVideo "go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLikeLogic {
	return &AddLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLikeLogic) AddLike(req *types.VideoLikeReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}
	// 调用videorpc判断视频是否存在
	_, err := l.svcCtx.VideoRpc.GetVideoById(l.ctx, &pbVideo.GetVideoByIdReq{Id: req.VideoId})
	if err != nil {
		return errors.Wrapf(err, "video_id: %d", req.VideoId)
	}
	// 获取当前用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 调用likerpc
	_, err = l.svcCtx.SocialRpc.VideoLike(l.ctx, &pb.VideoLikeReq{VideoId: req.VideoId, UserId: uid})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	return nil
}

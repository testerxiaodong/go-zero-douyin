package like

import (
	"context"
	"github.com/pkg/errors"
	pb2 "go-zero-douyin/apps/like/cmd/rpc/pb"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/like/cmd/api/internal/svc"
	"go-zero-douyin/apps/like/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLogic {
	return &AddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLogic) Add(req *types.VideoLikeReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := utils.GetValidator().ValidateZh(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}
	// 调用videorpc判断视频是否存在
	_, err := l.svcCtx.VideoRpc.GetVideoById(l.ctx, &pb.GetVideoByIdReq{Id: req.VideoId})
	if err != nil {
		return errors.Wrapf(err, "video_id: %d", req.VideoId)
	}
	// 获取当前用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 调用likerpc
	_, err = l.svcCtx.LikeRpc.VideoLike(l.ctx, &pb2.VideoLikeReq{VideoId: req.VideoId, UserId: uid})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	return nil
}

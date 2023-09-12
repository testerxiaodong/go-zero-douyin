package video

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelSectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelSectionLogic {
	return &DelSectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelSectionLogic) DelSection(req *types.DelSectionReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}

	// 调用videoRpc
	_, err := l.svcCtx.VideoRpc.DelSection(l.ctx, &pb.DelSectionReq{Id: req.Id})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}

	return nil
}

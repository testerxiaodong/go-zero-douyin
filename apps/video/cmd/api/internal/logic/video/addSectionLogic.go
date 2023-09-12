package video

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSectionLogic {
	return &AddSectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddSectionLogic) AddSection(req *types.AddSectionReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}

	// 调用videoRpc
	_, err := l.svcCtx.VideoRpc.AddSection(l.ctx, &pb.AddSectionReq{Name: req.Name})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}

	return nil
}

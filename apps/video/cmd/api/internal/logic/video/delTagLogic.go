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

type DelTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelTagLogic {
	return &DelTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelTagLogic) DelTag(req *types.DelTagReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}

	// 调用videorpc
	_, err := l.svcCtx.VideoRpc.DelTag(l.ctx, &pb.DelTagReq{Id: req.Id})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}

	return nil
}

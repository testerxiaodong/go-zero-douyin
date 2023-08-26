package user

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/user/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/user/cmd/api/internal/svc"
	"go-zero-douyin/apps/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateLogic {
	return &UpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateLogic) Update(req *types.UpdateUserReq) error {
	// todo: add your logic here and delete this line
	// 校验参数
	if validateResult := l.svcCtx.Validator.ValidateZh(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}
	// 获取用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)
	// rpc调用
	_, err := l.svcCtx.UserRpc.UpdateUser(l.ctx, &pb.UpdateUserReq{
		Id:       uid,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	return nil
}

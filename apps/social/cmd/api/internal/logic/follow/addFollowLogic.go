package follow

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

type AddFollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddFollowLogic {
	return &AddFollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddFollowLogic) AddFollow(req *types.UserFollowReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.ValidateZh(req); len(validateResult) > 0 {
		return errors.Wrapf(xerr.NewErrMsg(validateResult), "req: %v", req)
	}

	// 获取当前用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 调用social rpc
	_, err := l.svcCtx.SocialRpc.FollowUser(l.ctx, &pb.FollowUserReq{FollowerId: uid, UserId: req.UserId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	return nil
}

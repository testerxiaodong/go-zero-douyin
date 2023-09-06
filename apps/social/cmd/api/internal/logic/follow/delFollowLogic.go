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

type DelFollowLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDelFollowLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelFollowLogic {
	return &DelFollowLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DelFollowLogic) DelFollow(req *types.UserUnfollowReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return errors.Wrapf(xerr.NewErrMsg(validateResult), "req: %v", req)
	}

	// 获取当前用户uid
	uid := ctxdata.GetUidFromCtx(l.ctx)

	// 调用social rpc
	_, err := l.svcCtx.SocialRpc.UnfollowUser(l.ctx, &pb.UnfollowUserReq{FollowerId: uid, UserId: req.UserId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}

	return nil
}

package follow

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserFollowCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserFollowCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowCountLogic {
	return &GetUserFollowCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserFollowCountLogic) GetUserFollowCount(req *types.GetUserFollowCountReq) (resp *types.GetUserFollowCountResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 调用social rpc
	countResp, err := l.svcCtx.SocialRpc.GetUserFollowCount(l.ctx, &pb.GetUserFollowCountReq{UserId: req.UserId})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	resp = &types.GetUserFollowCountResp{}

	// 拷贝响应
	err = copier.Copy(resp, countResp)
	if err != nil {
		return nil, errors.Wrapf(err, "copy SocialRpc.GetUserFollowCount resp to api failed, data: %v", countResp)
	}

	return resp, nil
}

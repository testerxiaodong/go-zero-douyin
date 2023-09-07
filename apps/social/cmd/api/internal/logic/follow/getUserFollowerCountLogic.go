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

type GetUserFollowerCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserFollowerCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowerCountLogic {
	return &GetUserFollowerCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserFollowerCountLogic) GetUserFollowerCount(req *types.GetUserFollowerCountReq) (resp *types.GetUserFollowerCountResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 调用social rpc
	countResp, err := l.svcCtx.SocialRpc.GetUserFollowerCount(l.ctx, &pb.GetUserFollowerCountReq{UserId: req.UserId})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	// 拷贝响应
	resp = &types.GetUserFollowerCountResp{}
	err = copier.Copy(resp, countResp)
	if err != nil {
		return nil, errors.Wrapf(err, "copy SocialRpc.GetUserFollowerCount resp to api failed, data: %v", countResp)
	}

	return resp, nil
}

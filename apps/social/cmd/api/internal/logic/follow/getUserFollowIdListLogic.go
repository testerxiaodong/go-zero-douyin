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

type GetUserFollowIdListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserFollowIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowIdListLogic {
	return &GetUserFollowIdListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserFollowIdListLogic) GetUserFollowIdList(req *types.GetUserFollowIdListReq) (resp *types.GetUserFollowIdListResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg(validateResult), "req: %v", req)
	}

	// 调用 social rpc
	idListResp, err := l.svcCtx.SocialRpc.GetUserFollowIdList(l.ctx, &pb.GetUserFollowIdListReq{UserId: req.UserId})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	// 拷贝响应
	resp = &types.GetUserFollowIdListResp{UserIdList: make([]int64, 0)}
	err = copier.Copy(resp, idListResp)
	if err != nil {
		return nil, errors.Wrapf(err, "copy SocialRpc.GetUserFollowIdList resp to api failed, data: %v", idListResp)
	}

	return resp, nil
}

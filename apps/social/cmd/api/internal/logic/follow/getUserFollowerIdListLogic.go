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

type GetUserFollowerIdListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserFollowerIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserFollowerIdListLogic {
	return &GetUserFollowerIdListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserFollowerIdListLogic) GetUserFollowerIdList(req *types.GetUserFollowerIdListReq) (resp *types.GetUserFollowerIdListResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	// 调用social rpc
	idListResp, err := l.svcCtx.SocialRpc.GetUserFollowerIdList(l.ctx,
		&pb.GetUserFollowerIdListReq{UserId: req.UserId, Page: req.Page, PageSize: req.PageSize})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	// 拷贝响应
	resp = &types.GetUserFollowerIdListResp{UserIdList: make([]int64, 0)}
	err = copier.Copy(resp, idListResp)
	if err != nil {
		return nil, errors.Wrapf(err, "copy SocialRpc.GetUserFollowedIdList resp to api faield, data: %v", idListResp)
	}

	return resp, nil
}

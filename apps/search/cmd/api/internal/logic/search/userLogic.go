package search

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/search/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/search/cmd/api/internal/svc"
	"go-zero-douyin/apps/search/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) User(req *types.SearchUserReq) (resp *types.SearchUserResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}

	var sort pb.SearchUserReq_Sort
	if req.Sort == 1 {
		sort = pb.SearchUserReq_FollowerCOUNT
	} else {
		sort = pb.SearchUserReq_DEFAULT
	}

	// 调用searchRpc
	searchVideoResp, err := l.svcCtx.SearchRpc.SearchUser(l.ctx, &pb.SearchUserReq{
		Keyword:  req.Keyword,
		Page:     req.Page,
		PageSize: req.PageSize,
		Sort:     sort,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}

	// 拼接响应
	if len(searchVideoResp.GetUsers()) > 0 {
		resp = &types.SearchUserResp{Users: make([]*types.User, 0)}
		_ = copier.Copy(resp, searchVideoResp)
		return resp, nil
	}

	return
}

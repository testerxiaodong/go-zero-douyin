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

type CompleteVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteVideoLogic {
	return &CompleteVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteVideoLogic) CompleteVideo(req *types.CompleteVideoReq) (resp *types.CompleteVideoResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}
	// 调用searchRpc
	completeVideoResp, err := l.svcCtx.SearchRpc.CompleteVideo(l.ctx, &pb.CompleteVideoReq{Input: req.Input})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	// 拼接响应
	resp = &types.CompleteVideoResp{Suggestions: make([]string, 0)}
	_ = copier.Copy(resp, completeVideoResp)

	return resp, nil
}

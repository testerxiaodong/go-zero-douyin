package video

import (
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"

	"context"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteVideoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVideoLogic {
	return &DeleteVideoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteVideoLogic) DeleteVideo(req *types.DeleteVideoReq) error {
	// todo: add your logic here and delete this line
	if validateResult := utils.GetValidator().ValidateZh(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}
	return nil
}

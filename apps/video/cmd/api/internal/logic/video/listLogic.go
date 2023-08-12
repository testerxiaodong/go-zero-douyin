package video

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List(req *types.UserVideoListReq) (resp *types.UserVideoListResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := utils.GetValidator().ValidateZh(req); len(validateResult) > 0 {
		return nil, xerr.NewErrMsg(validateResult)
	}
	// 调用rpc
	videoList, err := l.svcCtx.VideoRpc.UserVideoList(l.ctx, &pb.UserVideoListReq{UserId: req.UserId})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	// 拼接响应
	resp = &types.UserVideoListResp{Videos: make([]*types.Video, 0)}
	err = copier.Copy(resp, videoList)
	if err != nil {
		return nil, errors.Wrapf(err, "copy data: %v", videoList)
	}
	fmt.Println(resp)
	return resp, nil
}

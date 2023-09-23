package video

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"go-zero-douyin/apps/video/cmd/api/internal/svc"
	"go-zero-douyin/apps/video/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllTagLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllTagLogic {
	return &GetAllTagLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllTagLogic) GetAllTag() (resp *types.GetAllTagResp, err error) {
	// todo: add your logic here and delete this line
	// 调用videoRpc
	tags, err := l.svcCtx.VideoRpc.GetAllTag(l.ctx, &pb.GetAllTagReq{})
	if err != nil {
		return nil, err
	}

	// 拷贝响应
	resp = &types.GetAllTagResp{Tags: make([]*types.Tag, 0)}
	err = copier.Copy(resp, tags)
	if err != nil {
		return nil, errors.Wrapf(err, "data: %v", tags)
	}

	return resp, nil
}

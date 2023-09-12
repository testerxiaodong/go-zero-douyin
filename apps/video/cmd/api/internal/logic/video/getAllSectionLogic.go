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

type GetAllSectionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllSectionLogic {
	return &GetAllSectionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllSectionLogic) GetAllSection() (resp *types.GetAllSectionResp, err error) {
	// todo: add your logic here and delete this line
	// 调用videoRpc
	sections, err := l.svcCtx.VideoRpc.GetAllSection(l.ctx, &pb.GetAllSectionReq{})
	if err != nil {
		return nil, err
	}

	// 拷贝响应
	resp = &types.GetAllSectionResp{Sections: make([]*types.SectionInfo, 0)}
	err = copier.Copy(resp, sections)
	if err != nil {
		return nil, errors.Wrapf(err, "data: %v", sections)
	}

	// 返回数据
	return resp, nil
}

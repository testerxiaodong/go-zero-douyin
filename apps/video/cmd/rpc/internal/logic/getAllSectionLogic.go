package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllSectionLogic {
	return &GetAllSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllSectionLogic) GetAllSection(in *pb.GetAllSectionReq) (*pb.GetAllSectionResp, error) {
	// todo: add your logic here and delete this line
	// 构建查询
	builder := l.svcCtx.SectionModel.SelectBuilder()
	sections, err := l.svcCtx.SectionModel.FindAll(l.ctx, builder, "")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询所有视频分区失败, err: %v", err)
	}
	// 长度为零
	if len(sections) == 0 {
		return &pb.GetAllSectionResp{}, nil
	}

	// 拷贝数据
	resp := &pb.GetAllSectionResp{Sections: make([]*pb.Section, 0)}
	_ = copier.Copy(&resp.Sections, sections)
	return resp, nil
}

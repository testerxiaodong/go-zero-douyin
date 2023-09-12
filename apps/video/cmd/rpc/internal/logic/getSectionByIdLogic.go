package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xerr"
	"gorm.io/gorm"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSectionByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSectionByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSectionByIdLogic {
	return &GetSectionByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSectionByIdLogic) GetSectionById(in *pb.GetSectionByIdReq) (*pb.GetSectionByIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get section by id with empty param")
	}

	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get section by id with empty id")
	}

	// 查询数据库
	section, err := l.svcCtx.SectionDo.GetSectionById(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"通过id从数据库中获取分区信息失败, err: %v, id: %d", err, in.GetId())
	}

	if section == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("分区不存在"), "section_id: %d", in.GetId())
	}

	// 返回数据
	return &pb.GetSectionByIdResp{
		SectionInfo: &pb.SectionInfo{
			Id:   section.ID,
			Name: section.Name}}, nil
}

package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddSectionLogic {
	return &AddSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddSection 添加分区
func (l *AddSectionLogic) AddSection(in *pb.AddSectionReq) (*pb.AddSectionResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add section with empty param")
	}
	if len(in.GetName()) == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add section with empty name")
	}

	// 查询分区名是否已存在
	section, err := l.svcCtx.SectionModel.FindOneByNameIsDelete(l.ctx, in.GetName(), xconst.DelStateNo)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "数据库根据名称查询分区失败, err: %v, name: %s", err, in.GetName())
	}

	if section != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("分区名已存在"), "新增分区失败， req: %v", in)
	}

	// 插入数据
	newSection := &model.Section{}
	newSection.Name = in.GetName()

	_, err = l.svcCtx.SectionModel.Insert(l.ctx, nil, newSection)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "数据库新增分区失败, err: %v, name: %s", err, in.GetName())
	}

	return &pb.AddSectionResp{}, nil
}

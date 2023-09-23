package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DelSectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelSectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelSectionLogic {
	return &DelSectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelSectionLogic) DelSection(in *pb.DelSectionReq) (*pb.DelSectionResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del section with empty param")
	}

	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del section with empty id")
	}
	// 查询分区是否存在
	section, err := l.svcCtx.SectionModel.FindOne(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询视频分区失败, section_id: %d", in.GetId())
	}
	if section == nil {
		return nil, errors.Wrapf(ErrSectionNotFound, "分区不存在, section_id: %d", in.GetId())
	}
	// 删除数据
	err = l.svcCtx.SectionModel.DeleteSoft(l.ctx, nil, section)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "删除分区失败, err: %v, id: %d", err, in.GetId())
	}

	return &pb.DelSectionResp{}, nil
}

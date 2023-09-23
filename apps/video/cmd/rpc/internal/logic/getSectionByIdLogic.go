package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrSectionNotFound = xerr.NewErrMsg("视频分区不存在")

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
	section, err := l.svcCtx.SectionModel.FindOne(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"通过id从数据库中获取分区信息失败, err: %v, id: %d", err, in.GetId())
	}

	if section == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("分区不存在"), "section_id: %d", in.GetId())
	}

	// 返回数据
	return &pb.GetSectionByIdResp{
		Section: &pb.Section{
			Id:   section.Id,
			Name: section.Name},
	}, nil
}

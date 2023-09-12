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

type GetTagByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTagByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTagByIdLogic {
	return &GetTagByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTagByIdLogic) GetTagById(in *pb.GetTagByIdReq) (*pb.GetTagByIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get tag by id with empty param")
	}

	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get tag by id with empty id")
	}

	// 查询数据库
	tag, err := l.svcCtx.TagDo.GetTagById(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR),
			"通过id从数据库中获取标签信息失败, err: %v, id: %d", err, in.GetId())
	}

	if tag == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("标签不存在"), "tag_id: %d", in.GetId())
	}

	// 返回数据
	return &pb.GetTagByIdResp{
		TagInfo: &pb.TagInfo{
			Id:   tag.ID,
			Name: tag.Name}}, nil
}

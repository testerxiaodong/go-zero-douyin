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

type DelTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDelTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DelTagLogic {
	return &DelTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DelTagLogic) DelTag(in *pb.DelTagReq) (*pb.DelTagResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del tag with empty param")
	}
	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "del tag with empty id")
	}
	// 查询视频标签是否存在
	tag, err := l.svcCtx.TagModel.FindOne(l.ctx, in.GetId())
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询视频标签失败, tag_id: %d", in.GetId())
	}
	if tag == nil {
		return nil, errors.Wrapf(ErrTagNotFound, "tag_id: %d", in.GetId())
	}
	// 删除数据
	err = l.svcCtx.TagModel.DeleteSoft(l.ctx, nil, tag)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_DELETE_ERR), "删除标签失败, err: %v, id: %d", err, in.GetId())
	}

	return &pb.DelTagResp{}, nil
}

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

type AddTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddTagLogic {
	return &AddTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddTag 新增标签
func (l *AddTagLogic) AddTag(in *pb.AddTagReq) (*pb.AddTagResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add tag with empty param")
	}

	if len(in.GetName()) == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "add tag with empty name")
	}

	// 查询标签名是否已存在
	tag, err := l.svcCtx.TagModel.FindOneByNameIsDelete(l.ctx, in.GetName(), xconst.DelStateNo)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "数据库根据名称查询标签失败, err: %v, name: %s", err, in.GetName())
	}

	if tag != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("标签名已存在"), "新增标签失败， req: %v", in)
	}

	// 插入数据
	newTag := &model.Tag{}
	newTag.Name = in.GetName()

	_, err = l.svcCtx.TagModel.Insert(l.ctx, nil, newTag)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "数据库新增标签失败, err: %v, name: %s", err, in.GetName())
	}

	return &pb.AddTagResp{}, nil
}

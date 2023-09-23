package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllTagLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllTagLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllTagLogic {
	return &GetAllTagLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllTagLogic) GetAllTag(in *pb.GetAllTagReq) (*pb.GetAllTagResp, error) {
	// todo: add your logic here and delete this line
	// 构建查询
	builder := l.svcCtx.TagModel.SelectBuilder()
	tags, err := l.svcCtx.TagModel.FindAll(l.ctx, builder, "")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "查询所有视频标签失败, err: %v", err)
	}
	// 长度为零
	if len(tags) == 0 {
		return &pb.GetAllTagResp{}, nil
	}

	// 拼接数据
	resp := &pb.GetAllTagResp{Tags: make([]*pb.Tag, 0, len(tags))}
	err = copier.Copy(&resp.Tags, tags)
	if err != nil {
		l.Logger.Errorf("拷贝所有标签数据失败, err: %v", err)
	}
	return resp, nil
}

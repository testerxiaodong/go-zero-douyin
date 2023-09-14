package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/search/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/search/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVideoLogic {
	return &DeleteVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteVideoLogic) DeleteVideo(in *pb.DeleteVideoDocumentReq) (*pb.DeleteVideoDocumentResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "参数不能为nil")
	}
	if in.GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "视频id不能为空")
	}
	// 调用es
	_, err := l.svcCtx.ElasticSearch.DeleteDocument(l.ctx, xconst.ElasticSearchVideoIndexName, cast.ToString(in.GetId()))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_DELETE_ERR), "删除视频文档失败, err: %v", err)
	}
	return &pb.DeleteVideoDocumentResp{}, nil
}

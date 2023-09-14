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

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserLogic) DeleteUser(in *pb.DeleteUserDocumentReq) (*pb.DeleteUserDocumentResp, error) {
	// todo: add your logic here and delete this line
	// 调用es
	_, err := l.svcCtx.ElasticSearch.DeleteDocument(l.ctx, xconst.ElasticSearchUserIndexName, cast.ToString(in.GetId()))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_DELETE_ERR), "删除用户文档失败, err: %v", err)
	}
	return &pb.DeleteUserDocumentResp{}, nil
}

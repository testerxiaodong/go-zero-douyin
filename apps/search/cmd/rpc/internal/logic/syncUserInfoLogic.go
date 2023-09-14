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

type SyncUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncUserInfoLogic {
	return &SyncUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SyncUserInfo 创建/更新用户信息
func (l *SyncUserInfoLogic) SyncUserInfo(in *pb.SyncUserInfoReq) (*pb.SyncUserInfoResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in.GetUser().GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "用户id不允许为空")
	}
	// 调用es接口
	_, err := l.svcCtx.ElasticSearch.CreateDocument(l.ctx, xconst.ElasticSearchUserIndexName, cast.ToString(in.GetUser().GetId()), in.GetUser())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "更新es文档失败, err: %v", err)
	}
	return &pb.SyncUserInfoResp{}, nil
}

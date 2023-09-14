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

type SyncVideoInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncVideoInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncVideoInfoLogic {
	return &SyncVideoInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncVideoInfoLogic) SyncVideoInfo(in *pb.SyncVideoInfoReq) (*pb.SyncVideoInfoResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil || in.GetVideo() == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "参数不能为nil")
	}
	if in.GetVideo().GetId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "视频id不允许为空")
	}
	// 调用es接口
	_, err := l.svcCtx.ElasticSearch.CreateDocument(l.ctx, xconst.ElasticSearchVideoIndexName, cast.ToString(in.GetVideo().GetId()), in.GetVideo())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "更新es视频文档失败, err: %v", err)
	}
	return &pb.SyncVideoInfoResp{}, nil
}

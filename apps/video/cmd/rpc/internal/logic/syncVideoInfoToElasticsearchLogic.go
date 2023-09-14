package logic

import (
	"context"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncVideoInfoToElasticsearchLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncVideoInfoToElasticsearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncVideoInfoToElasticsearchLogic {
	return &SyncVideoInfoToElasticsearchLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncVideoInfoToElasticsearchLogic) SyncVideoInfoToElasticsearch(in *pb.SyncVideoInfoToElasticsearchReq) (*pb.SyncVideoInfoToElasticsearchResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SyncVideoInfoToElasticsearchResp{}, nil
}

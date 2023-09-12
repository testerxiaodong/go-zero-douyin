package logic

import (
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go-zero-douyin/common/xerr"

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
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "同步视频信息到es的参数为空")
	}

	// 全量更新
	_, err := l.svcCtx.Elasticsearch.CreateDocument(l.ctx, "video", cast.ToString(in.GetVideoInfo().GetId()), in.GetVideoInfo())
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_INSERT_ERR), "创建es文档失败, err: %v, video: %v", err, in.GetVideoInfo())
	}
	return &pb.SyncVideoInfoToElasticsearchResp{}, nil
}

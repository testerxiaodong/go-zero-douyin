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
	// 查询数据库所有视频信息数据
	videos, err := l.svcCtx.VideoDo.GetAllVideo(l.ctx)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "从数据库中获取所有视频信息失败, err: %v", err)
	}

	// 没有数据，直接返回
	if len(videos) == 0 {
		return &pb.SyncVideoInfoToElasticsearchResp{}, nil
	}

	// 全量更新
	for _, video := range videos {
		idStr := cast.ToString(video.ID)
		_, err := l.svcCtx.Elasticsearch.CreateDocument(l.ctx, "video", idStr, video)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_INSERT_ERR), "创建es文档失败, err: %v, video: %v", err, video)
		}
	}
	return &pb.SyncVideoInfoToElasticsearchResp{}, nil
}

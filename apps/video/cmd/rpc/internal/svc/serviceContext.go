package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-zero-douyin/apps/video/cmd/rpc/internal/config"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	taskQueue "go-zero-douyin/common/asynq"
)

type ServiceContext struct {
	Config       config.Config
	VideoModel   model.VideoModel
	TagModel     model.TagModel
	SectionModel model.SectionModel
	Asynq        taskQueue.TaskQueueClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config:       c,
		VideoModel:   model.NewVideoModel(sqlConn, c.Cache),
		TagModel:     model.NewTagModel(sqlConn, c.Cache),
		SectionModel: model.NewSectionModel(sqlConn, c.Cache),
		Asynq:        taskQueue.NewAsynq(asynq.NewClient(asynq.RedisClientOpt{Addr: c.RedisConf.Host})),
	}
}

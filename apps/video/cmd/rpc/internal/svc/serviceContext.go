package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/syncx"
	"go-zero-douyin/apps/video/cmd/rpc/internal/config"
	"go-zero-douyin/apps/video/cmd/rpc/internal/dao"
	taskQueue "go-zero-douyin/common/asynq"
	"go-zero-douyin/common/cache"
)

type ServiceContext struct {
	Config       config.Config
	VideoDo      dao.VideoDo
	SectionDo    dao.SectionDo
	TagDo        dao.TagDo
	Redis        cache.RedisCache
	Rabbit       rabbitmq.Sender
	Asynq        taskQueue.TaskQueueClient
	SingleFlight syncx.SingleFlight
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		VideoDo:      dao.NewVideoRepository(c.DataSource),
		SectionDo:    dao.NewSectionRepository(c.DataSource),
		TagDo:        dao.NewTagRepository(c.DataSource),
		Redis:        cache.NewRedisClient(c.RedisCache),
		Rabbit:       rabbitmq.MustNewSender(c.RabbitSenderConf),
		Asynq:        taskQueue.NewAsynq(asynq.NewClient(asynq.RedisClientOpt{Addr: c.Redis.Host})),
		SingleFlight: syncx.NewSingleFlight(),
	}
}

package svc

import (
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/syncx"
	"go-zero-douyin/apps/social/cmd/rpc/internal/config"
	"go-zero-douyin/apps/social/cmd/rpc/internal/dao"
	"go-zero-douyin/common/cache"
)

type ServiceContext struct {
	Config       config.Config
	CommentDo    dao.CommentDo
	FollowDo     dao.FollowDo
	LikeDo       dao.LikeDo
	Rabbit       rabbitmq.Sender
	Redis        cache.RedisCache
	SingleFlight syncx.SingleFlight
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		CommentDo:    dao.NewCommentRepository(c.DataSource),
		FollowDo:     dao.NewFollowRepository(c.DataSource),
		LikeDo:       dao.NewLikeRepository(c.DataSource),
		Rabbit:       rabbitmq.MustNewSender(c.RabbitSenderConf),
		Redis:        cache.NewRedisClient(c.RedisCache),
		SingleFlight: syncx.NewSingleFlight(),
	}
}

package svc

import (
	"github.com/zeromicro/go-queue/rabbitmq"
	"go-zero-douyin/apps/user/cmd/rpc/internal/config"
	"go-zero-douyin/apps/user/cmd/rpc/internal/dao"
)

type ServiceContext struct {
	Config config.Config
	UserDo dao.UserDo
	Rabbit rabbitmq.Sender
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		UserDo: dao.NewUserRepository(c.DataSource),
		Rabbit: rabbitmq.MustNewSender(c.RabbitSenderConf),
	}
}

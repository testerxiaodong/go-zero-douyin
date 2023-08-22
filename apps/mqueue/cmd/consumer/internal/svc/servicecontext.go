package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.Redis),
	}
}

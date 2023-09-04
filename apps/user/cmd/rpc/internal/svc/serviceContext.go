package svc

import (
	"go-zero-douyin/apps/user/cmd/rpc/internal/config"
	"go-zero-douyin/apps/user/cmd/rpc/internal/dao"
)

type ServiceContext struct {
	Config config.Config
	UserDo dao.UserDo
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		UserDo: dao.NewUserRepository(c.DataSource),
	}
}

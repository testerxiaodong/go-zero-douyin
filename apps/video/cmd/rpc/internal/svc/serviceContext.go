package svc

import (
	"go-zero-douyin/apps/video/cmd/rpc/internal/config"
	"go-zero-douyin/apps/video/cmd/rpc/internal/dao"
)

type ServiceContext struct {
	Config  config.Config
	VideoDo dao.VideoDo
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		VideoDo: dao.NewVideoRepository(c.DataSource),
	}
}

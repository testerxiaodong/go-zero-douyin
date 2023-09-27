package svc

import (
	"go-zero-douyin/apps/recommend/cmd/rpc/internal/config"
	"go-zero-douyin/common/http"
)

type ServiceContext struct {
	Config     config.Config
	GorseConf  string
	RestClient http.RestClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		GorseConf:  c.GorseConf,
		RestClient: http.NewRestyClient(),
	}
}

package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/user/cmd/api/internal/config"
	"go-zero-douyin/apps/user/cmd/rpc/userrpc"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc userrpc.Userrpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: userrpc.NewUserrpc(zrpc.MustNewClient(c.UserRpcConf)),
	}
}

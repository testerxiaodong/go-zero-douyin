package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/config"
	"go-zero-douyin/apps/social/cmd/rpc/social"
)

type ServiceContext struct {
	Config    config.Config
	SocialRpc social.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		SocialRpc: social.NewSocial(zrpc.MustNewClient(c.SocialRpcConf)),
	}
}

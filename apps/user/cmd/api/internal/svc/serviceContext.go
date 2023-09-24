package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/social/cmd/rpc/social"
	"go-zero-douyin/apps/user/cmd/api/internal/config"
	"go-zero-douyin/apps/user/cmd/rpc/user"
	"go-zero-douyin/common/utils"
)

type ServiceContext struct {
	Config    config.Config
	UserRpc   user.User
	SocialRpc social.Social
	Validator utils.Validator
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		UserRpc:   user.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
		SocialRpc: social.NewSocial(zrpc.MustNewClient(c.SocialRpcConf)),
		Validator: utils.NewZhValidator(),
	}
}

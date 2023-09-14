package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/search/cmd/api/internal/config"
	"go-zero-douyin/apps/search/cmd/rpc/search"
	"go-zero-douyin/common/utils"
)

type ServiceContext struct {
	Config    config.Config
	SearchRpc search.Search
	Validator utils.Validator
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		SearchRpc: search.NewSearch(zrpc.MustNewClient(c.SearchRpcConf)),
		Validator: utils.NewZhValidator(),
	}
}

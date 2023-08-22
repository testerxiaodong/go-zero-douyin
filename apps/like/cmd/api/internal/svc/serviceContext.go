package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/like/cmd/api/internal/config"
	"go-zero-douyin/apps/like/cmd/rpc/likerpc"
	"go-zero-douyin/apps/video/cmd/rpc/videorpc"
)

type ServiceContext struct {
	Config   config.Config
	LikeRpc  likerpc.Likerpc
	VideoRpc videorpc.Videorpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:   c,
		LikeRpc:  likerpc.NewLikerpc(zrpc.MustNewClient(c.LikeRpcConf)),
		VideoRpc: videorpc.NewVideorpc(zrpc.MustNewClient(c.VideoRpcConf)),
	}
}

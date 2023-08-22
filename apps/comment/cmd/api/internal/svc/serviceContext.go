package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/comment/cmd/api/internal/config"
	"go-zero-douyin/apps/comment/cmd/rpc/commentrpc"
	"go-zero-douyin/apps/video/cmd/rpc/videorpc"
)

type ServiceContext struct {
	Config     config.Config
	CommentRpc commentrpc.Commentrpc
	VideoRpc   videorpc.Videorpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:     c,
		CommentRpc: commentrpc.NewCommentrpc(zrpc.MustNewClient(c.CommentRpcConf)),
		VideoRpc:   videorpc.NewVideorpc(zrpc.MustNewClient(c.VideoRpcConf)),
	}
}

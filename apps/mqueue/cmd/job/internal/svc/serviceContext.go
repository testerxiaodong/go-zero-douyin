package svc

import (
	"github.com/hibiken/asynq"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/mqueue/cmd/job/internal/config"
	"go-zero-douyin/apps/video/cmd/rpc/video"
)

type ServiceContext struct {
	Config      config.Config
	AsynqServer *asynq.Server
	VideoRpc    video.Video
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		AsynqServer: newAsynqServer(c),
		VideoRpc:    video.NewVideo(zrpc.MustNewClient(c.VideoRpcConf)),
	}
}

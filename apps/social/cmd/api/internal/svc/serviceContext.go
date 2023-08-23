package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/social/cmd/api/internal/config"
	"go-zero-douyin/apps/social/cmd/rpc/social"
	"go-zero-douyin/apps/video/cmd/rpc/video"
)

type ServiceContext struct {
	Config    config.Config
	SocialRpc social.Social
	VideoRpc  video.Video
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		SocialRpc: social.NewSocial(zrpc.MustNewClient(c.SocialRpcConf)),
		VideoRpc:  video.NewVideo(zrpc.MustNewClient(c.VideoRpcConf)),
	}
}

package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/recommend/cmd/api/internal/config"
	"go-zero-douyin/apps/recommend/cmd/rpc/recommend"
	"go-zero-douyin/apps/social/cmd/rpc/social"
	"go-zero-douyin/apps/video/cmd/rpc/video"
	"go-zero-douyin/common/utils"
)

type ServiceContext struct {
	Config       config.Config
	VideoRpc     video.Video
	SocialRpc    social.Social
	RecommendRpc recommend.Recommend
	Validator    utils.Validator
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		VideoRpc:     video.NewVideo(zrpc.MustNewClient(c.VideoRpcConf)),
		SocialRpc:    social.NewSocial(zrpc.MustNewClient(c.SocialRpcConf)),
		RecommendRpc: recommend.NewRecommend(zrpc.MustNewClient(c.RecommendRpcConf)),
		Validator:    utils.NewZhValidator(),
	}
}

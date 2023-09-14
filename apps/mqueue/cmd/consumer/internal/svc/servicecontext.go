package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/config"
	"go-zero-douyin/apps/search/cmd/rpc/search"
	"go-zero-douyin/apps/social/cmd/rpc/social"
	"go-zero-douyin/apps/user/cmd/rpc/user"
	"go-zero-douyin/apps/video/cmd/rpc/video"
)

type ServiceContext struct {
	Config    config.Config
	Redis     *redis.Redis
	UserRpc   user.User
	VideoRpc  video.Video
	SocialRpc social.Social
	SearchRpc search.Search
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		Redis:     redis.MustNewRedis(c.Redis),
		UserRpc:   user.NewUser(zrpc.MustNewClient(c.UserRpcConf)),
		VideoRpc:  video.NewVideo(zrpc.MustNewClient(c.VideoRpcConf)),
		SocialRpc: social.NewSocial(zrpc.MustNewClient(c.SocialRpcConf)),
		SearchRpc: search.NewSearch(zrpc.MustNewClient(c.SearchRpcConf)),
	}
}

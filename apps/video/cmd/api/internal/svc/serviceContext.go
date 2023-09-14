package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/search/cmd/rpc/search"
	"go-zero-douyin/apps/social/cmd/rpc/social"
	"go-zero-douyin/apps/video/cmd/api/internal/config"
	"go-zero-douyin/apps/video/cmd/rpc/video"
	"go-zero-douyin/common/utils"
)

type ServiceContext struct {
	Config    config.Config
	VideoRpc  video.Video
	SocialRpc social.Social
	SearchRpc search.Search
	Validator utils.Validator
	OssClient utils.OssClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		VideoRpc:  video.NewVideo(zrpc.MustNewClient(c.VideoRpcConf)),
		SocialRpc: social.NewSocial(zrpc.MustNewClient(c.SocialRpcConf)),
		SearchRpc: search.NewSearch(zrpc.MustNewClient(c.SearchRpcConf)),
		Validator: utils.NewZhValidator(),
		OssClient: utils.NewAliOssClient(c.AliCloud.AccessKeyId, c.AliCloud.AccessKeySecret, c.AliCloud.EndPoint, c.AliCloud.BucketName),
	}
}

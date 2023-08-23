package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/social/cmd/rpc/social"
	"go-zero-douyin/apps/video/cmd/api/internal/config"
	"go-zero-douyin/apps/video/cmd/rpc/video"
	"go-zero-douyin/common/utils"
)

type ServiceContext struct {
	Config    config.Config
	VideoRpc  video.Video
	SocialRpc social.Social
	Validator *utils.Validator
	OssClient *utils.OssClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		VideoRpc:  video.NewVideo(zrpc.MustNewClient(c.VideoRpcConf)),
		SocialRpc: social.NewSocial(zrpc.MustNewClient(c.SocialRpcConf)),
		Validator: utils.GetValidator(),
		OssClient: utils.InitOssClient(c.AliCloud.AccessKeyId, c.AliCloud.AccessKeySecret, c.AliCloud.EndPoint, c.AliCloud.BucketName),
	}
}

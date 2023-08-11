package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-zero-douyin/apps/video/cmd/api/internal/config"
	"go-zero-douyin/apps/video/cmd/rpc/videorpc"
	"go-zero-douyin/common/utils"
)

type ServiceContext struct {
	Config    config.Config
	VideoRpc  videorpc.Videorpc
	Validator *utils.Validator
	OssClient *utils.OssClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		VideoRpc:  videorpc.NewVideorpc(zrpc.MustNewClient(c.VideoRpcConf)),
		Validator: utils.GetValidator(),
		OssClient: utils.InitOssClient(c.AliCloud.AccessKeyId, c.AliCloud.AccessKeySecret, c.AliCloud.EndPoint, c.AliCloud.BucketName),
	}
}

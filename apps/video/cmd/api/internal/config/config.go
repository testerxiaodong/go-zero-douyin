package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UserRpcConf   zrpc.RpcClientConf
	VideoRpcConf  zrpc.RpcClientConf
	SocialRpcConf zrpc.RpcClientConf
	SearchRpcConf zrpc.RpcClientConf
	JwtAuth       struct {
		AccessSecret string
	}
	AliCloud struct {
		AccessKeyId     string
		AccessKeySecret string
		BucketName      string
		EndPoint        string
		CommonPath      string
	}
}

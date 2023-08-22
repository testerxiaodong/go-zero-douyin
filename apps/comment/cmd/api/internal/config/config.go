package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	CommentRpcConf zrpc.RpcClientConf
	VideoRpcConf   zrpc.RpcClientConf
	JwtAuth        struct {
		AccessSecret string
	}
}

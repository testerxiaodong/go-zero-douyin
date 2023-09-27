package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	// kqs : pub sub
	VideoLikeConf kq.KqConf
	// rpc
	SocialRpcConf zrpc.RpcClientConf
}

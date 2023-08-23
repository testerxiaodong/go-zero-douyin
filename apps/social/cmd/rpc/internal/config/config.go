package config

import (
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DataSource       string
	RabbitSenderConf rabbitmq.RabbitSenderConf
	RedisCache       redis.RedisConf
}
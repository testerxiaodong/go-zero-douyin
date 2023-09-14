package config

import (
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf
	Redis                          redis.RedisConf
	UserRpcConf                    zrpc.RpcClientConf
	VideoRpcConf                   zrpc.RpcClientConf
	SocialRpcConf                  zrpc.RpcClientConf
	SearchRpcConf                  zrpc.RpcClientConf
	RabbitVideoSectionMqConf       rabbitmq.RabbitListenerConf
	RabbitVideoTagMqConf           rabbitmq.RabbitListenerConf
	RabbitVideoCommentMqConf       rabbitmq.RabbitListenerConf
	RabbitUserLikeVideoMqConf      rabbitmq.RabbitListenerConf
	RabbitVideoLikedByUserMqConf   rabbitmq.RabbitListenerConf
	RabbitUserFollowUserMqConf     rabbitmq.RabbitListenerConf
	RabbitUserFollowedByUserMqConf rabbitmq.RabbitListenerConf
	RabbitMysqlUserDeleteMqConf    rabbitmq.RabbitListenerConf
	RabbitMysqlUserUpdateMqConf    rabbitmq.RabbitListenerConf
	RabbitMysqlVideoDeleteMqConf   rabbitmq.RabbitListenerConf
	RabbitMysqlVideoUpdateMqConf   rabbitmq.RabbitListenerConf
}

package config

import (
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	service.ServiceConf
	Redis                          redis.RedisConf
	RabbitVideoCommentMqConf       rabbitmq.RabbitListenerConf
	RabbitUserLikeVideoMqConf      rabbitmq.RabbitListenerConf
	RabbitVideoLikedByUserMqConf   rabbitmq.RabbitListenerConf
	RabbitUserFollowUserMqConf     rabbitmq.RabbitListenerConf
	RabbitUserFollowedByUserMqConf rabbitmq.RabbitListenerConf
}

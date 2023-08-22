package listen

import (
	"context"
	"github.com/zeromicro/go-zero/core/service"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/config"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
)

// Mqs 所有的消息队列
func Mqs(c config.Config) []service.Service {

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()

	var services []service.Service

	// RabbitMq ：pub sub
	services = append(services, RabbitMqs(c, ctx, svcContext)...)

	return services
}

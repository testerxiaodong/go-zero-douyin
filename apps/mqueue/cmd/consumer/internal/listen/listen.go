package listen

import (
	"context"
	"github.com/zeromicro/go-zero/core/service"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/config"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
)

// Mqs back to all consumers
func Mqs(c config.Config) []service.Service {

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()

	var services []service.Service

	//kqs ：pub sub
	services = append(services, KqMqs(c, ctx, svcContext)...)

	return services
}

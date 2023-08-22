package listen

import (
	"context"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/service"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/config"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/mqs/rmqs"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
)

// RabbitMqs rabbitmq消息队列
func RabbitMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		// 视频评论数的缓存一致性消息
		rabbitmq.MustNewListener(c.RabbitListenerConf, rmqs.NewVideoCommentMq(ctx, svcContext)),
		// 用户点赞视频
		rabbitmq.MustNewListener(c.RabbitListenerConf, rmqs.NewUserLikeVideoMq(ctx, svcContext)),
		// 视频被用户点赞
		rabbitmq.MustNewListener(c.RabbitListenerConf, rmqs.NewVideoLikedByUserMq(ctx, svcContext)),
	}

}

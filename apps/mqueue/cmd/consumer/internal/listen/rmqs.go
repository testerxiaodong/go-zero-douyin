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
		// 视频分区
		rabbitmq.MustNewListener(c.RabbitVideoSectionMqConf, rmqs.NewVideoSectionMq(ctx, svcContext)),
		// 视频标签
		rabbitmq.MustNewListener(c.RabbitVideoTagMqConf, rmqs.NewVideoTagMq(ctx, svcContext)),
		// 视频评论数的缓存一致性消息
		rabbitmq.MustNewListener(c.RabbitVideoCommentMqConf, rmqs.NewVideoCommentMq(ctx, svcContext)),
		// 用户点赞视频
		rabbitmq.MustNewListener(c.RabbitUserLikeVideoMqConf, rmqs.NewUserLikeVideoMq(ctx, svcContext)),
		// 视频被用户点赞
		rabbitmq.MustNewListener(c.RabbitVideoLikedByUserMqConf, rmqs.NewVideoLikedByUserMq(ctx, svcContext)),
		// 用户关注
		rabbitmq.MustNewListener(c.RabbitUserFollowUserMqConf, rmqs.NewUserFollowUserMq(ctx, svcContext)),
		// 用户被关注
		rabbitmq.MustNewListener(c.RabbitUserFollowedByUserMqConf, rmqs.NewUserFollowedByUserMq(ctx, svcContext)),
		// 删除用户es文档
		rabbitmq.MustNewListener(c.RabbitMysqlUserDeleteMqConf, rmqs.NewMysqlUserDeleteMq(ctx, svcContext)),
		// 更新用户es文档
		rabbitmq.MustNewListener(c.RabbitMysqlUserUpdateMqConf, rmqs.NewMysqlUserUpdateMq(ctx, svcContext)),
		// 删除视频es文档
		rabbitmq.MustNewListener(c.RabbitMysqlVideoDeleteMqConf, rmqs.NewMysqlVideoDeleteMq(ctx, svcContext)),
		// 更新视频es文档
		rabbitmq.MustNewListener(c.RabbitMysqlVideoUpdateMqConf, rmqs.NewMysqlVideoUpdateMq(ctx, svcContext)),
	}

}

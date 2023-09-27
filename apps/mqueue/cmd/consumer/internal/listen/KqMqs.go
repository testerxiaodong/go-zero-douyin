package listen

import (
	"context"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/config"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/mqs/kqs"
	"go-zero-douyin/apps/mqueue/cmd/consumer/internal/svc"
)

// KqMqs pub sub use kqs (kafka)
func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		// 视频点赞消息
		kq.MustNewQueue(c.VideoLikeConf, kqs.NewVideoLikeMq(ctx, svcContext)),
	}
}

package logic

import (
	"context"
	"github.com/hibiken/asynq"
	"go-zero-douyin/apps/mqueue/cmd/job/internal/svc"
	"go-zero-douyin/apps/mqueue/cmd/job/jobtype"
)

type CronJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext) *CronJob {
	return &CronJob{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Register 注册任务路由
func (l *CronJob) Register() *asynq.ServeMux {

	mux := asynq.NewServeMux()

	//defer job
	mux.Handle(jobtype.DeferPublishVideo, NewPublishVideoHandler(l.svcCtx))

	//queue job , asynq support queue job

	return mux
}

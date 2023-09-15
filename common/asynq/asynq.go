package asynq

import (
	"context"
	"github.com/hibiken/asynq"
)

type TaskQueueClient interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
	EnqueueContext(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type Asynq struct {
	client *asynq.Client
}

func NewAsynq(client *asynq.Client) TaskQueueClient {
	return &Asynq{
		client: client,
	}
}

func (q *Asynq) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return q.client.Enqueue(task, opts...)
}

func (q *Asynq) EnqueueContext(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return q.client.EnqueueContext(ctx, task, opts...)
}

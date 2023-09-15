package svc

import (
	"fmt"
	"github.com/hibiken/asynq"
	"go-zero-douyin/apps/mqueue/cmd/job/internal/config"
)

func newAsynqServer(c config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Redis.Host},
		asynq.Config{
			IsFailure: func(err error) bool {
				fmt.Printf("asynq server exec task IsFailure ======== >>>>>>>>>>>  err : %+v \n", err)
				return true
			},
			Concurrency: 20,
		},
	)
}

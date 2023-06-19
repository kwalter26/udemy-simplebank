package worker

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
)

const (
	EmailQueue   = "email"
	DefaultQueue = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			EmailQueue:   10,
			DefaultQueue: 5,
		},
	})
	return &RedisTaskProcessor{server: server, store: store}
}
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}

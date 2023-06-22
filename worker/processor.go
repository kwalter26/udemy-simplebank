package worker

import (
	"context"
	"github.com/hibiken/asynq"
	db "github.com/kwalter26/udemy-simplebank/db/sqlc"
	"github.com/kwalter26/udemy-simplebank/mail"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, mailer mail.EmailSender) TaskProcessor {
	logger := NewLogger()
	redis.SetLogger(logger)
	server := asynq.NewServer(redisOpt, asynq.Config{

		Concurrency: 10,
		Queues: map[string]int{
			EmailQueue:   10,
			DefaultQueue: 5,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			// log error messages to stdout
			log.Error().
				Err(err).
				Str("type", task.Type()).
				Bytes("payload", task.Payload()).
				Msg("process task error")
		}),
		Logger: logger,
	})
	return &RedisTaskProcessor{server: server, store: store, mailer: mailer}
}
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}

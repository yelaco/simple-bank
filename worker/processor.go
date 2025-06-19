// Package worker provides functionality for processing asynchronous tasks
// using Redis as the message broker.
package worker

import (
	"context"

	"github.com/hibiken/asynq"
	db "github.com/yelaco/simple-bank/db/sqlc"
)

// TaskProcessor defines the interface for processing asynchronous tasks.
type TaskProcessor interface {
	// Start initializes the task processor and begins processing tasks.
	Start() error

	// ProcessTaskSendVerifyEmail processes the "send verify email" task.
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

// RedisTaskProcessor implements the TaskProcessor interface using Redis
// as the message broker for asynchronous task processing.
type RedisTaskProcessor struct {
	server *asynq.Server // Redis server instance for task processing
	store  db.Store      // Database store for data operations
}

// Start initializes the Redis task processor and begins listening for tasks.
// It registers the task handlers and starts the Redis server.
func (processor *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	processor.server.Start(mux)

	return nil
}

// NewRedisTaskProcessor creates a new instance of RedisTaskProcessor
// with the provided Redis client options and database store.
func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{})

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}

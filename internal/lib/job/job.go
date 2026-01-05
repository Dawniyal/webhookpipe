package job

import (
	"context"
	"encoding/json"

	"github.com/Dawniyal/webhookpipe/internal/config"
	"github.com/Dawniyal/webhookpipe/internal/lib/forward"
	"github.com/Dawniyal/webhookpipe/internal/model/event"
	"github.com/Dawniyal/webhookpipe/internal/repository"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type JobService struct {
	Client    *asynq.Client
	server    *asynq.Server
	logger    *zerolog.Logger
	forward   *forward.Forward
	eventRepo *repository.EventsRepository
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config, eventRepo *repository.EventsRepository) *JobService {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       0,
	})

	failedHandler := asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
		var p ForwardPayload
		json.Unmarshal(task.Payload(), &p)

		eventRepo.UpdateStatus(ctx, p.EventID, event.StatusFailed)
	})

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.Redis.Address, Password: cfg.Redis.Password, DB: 0},
		asynq.Config{
			Concurrency:  25,
			ErrorHandler: failedHandler,
		},
	)

	return &JobService{
		Client:    client,
		server:    server,
		logger:    logger,
		forward:   forward.NewForward(logger),
		eventRepo: eventRepo,
	}
}

func (j *JobService) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSend, j.handleForward)

	j.logger.Info().Msg("Starting asynq")

	if err := j.server.Start(mux); err != nil {
		return err
	}

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping asynq")
	j.server.Shutdown()
	j.Client.Close()
}

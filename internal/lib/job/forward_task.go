package job

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Dawniyal/webhookpipe/internal/model/event"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	TaskSend = "forward:send"
)

var validate = validator.New()

type ForwardPayload struct {
	EventID   uuid.UUID      `json:"eventId"`
	TargetURL string         `json:"targetUrl"`
	Payload   map[string]any `json:"payload"`
}

func NewForward(eventID uuid.UUID, targetURL string, payload map[string]any) (*asynq.Task, error) {
	p, err := json.Marshal(ForwardPayload{
		EventID:   eventID,
		TargetURL: targetURL,
		Payload:   payload,
	})

	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TaskSend, p,
		asynq.MaxRetry(5),
		asynq.Timeout(10*time.Second)), nil
}

func (j *JobService) handleForward(ctx context.Context, t *asynq.Task) error {
	var p ForwardPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	j.logger.Info().
		Str("type", "forward").
		Str("target_url", p.TargetURL).
		Msg("processing forward Task")

	err := validate.Var(p.TargetURL, "http_url")
	if err != nil {
		j.logger.Error().Err(err).Msgf("failed to validate target_url: %v", p.TargetURL)
		return err
	}

	err = j.forward.ForwardPayload(
		p.TargetURL,
		p.Payload,
	)
	if err != nil {
		j.logger.Error().Err(err).Msg("failed to forward payload")
		return err
	}

	err = j.eventRepo.UpdateStatus(ctx, p.EventID, event.StatusDelivered)
	if err != nil {
		j.logger.Error().Err(err).Msg("failed to change status of event")
	}

	j.logger.Info().
		Str("type", "forward").
		Str("target_url", p.TargetURL).
		Msg("successfully forward Task")

	return nil

}

func (j *JobService) EnqueueForward(eventID uuid.UUID, targetURL string, payload map[string]any) error {
	task, err := NewForward(eventID, targetURL, payload)
	if err != nil {
		return err
	}

	info, err := j.Client.Enqueue(task)
	if err != nil {
		j.logger.Error().Err(err).Msg("failed to enqueue task")
		return err
	}

	j.logger.Info().Str("task_id", info.ID).Msg("enqueued task")
	return nil
}

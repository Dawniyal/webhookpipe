package service

import (
	"context"

	"github.com/Dawniyal/webhookpipe/internal/lib/job"
	"github.com/Dawniyal/webhookpipe/internal/middleware"
	"github.com/Dawniyal/webhookpipe/internal/model/event"
	"github.com/Dawniyal/webhookpipe/internal/repository"
)

type EventService struct {
	job       *job.JobService
	eventRepo *repository.EventsRepository
}

func NewEventService(job *job.JobService, eventRepo *repository.EventsRepository) *EventService {
	return &EventService{
		job:       job,
		eventRepo: eventRepo,
	}
}

func (s *EventService) CreateEvent(
	ctx context.Context,
	payload *event.CreateEventPayload,
) (*event.Event, error) {
	logger := middleware.GetLogger(ctx)

	logger.Info().
		Str("endpoint_id", payload.EndpointID).
		Msg("creating event")

	evt, err := s.eventRepo.CreateEvent(ctx, payload)
	if err != nil {
		logger.Error().
			Err(err).
			Str("endpoint_id", payload.EndpointID).
			Msg("failed to create event")
		return nil, err
	}

	err = s.job.EnqueueForward(evt.ID, evt.EndpointID, evt.Payload)

	if err != nil {
		logger.Error().
			Err(err).
			Str("event_id", evt.ID.String()).
			Msg("failed to enqueue forward job")
	}

	err = s.eventRepo.UpdateStatus(ctx, evt.ID, event.StatusPending)
	if err != nil {
		logger.Error().
			Err(err).
			Str("event_id", evt.ID.String()).
			Msg("failed to update event status")
		return nil, err
	}

	logger.Info().
		Str("event_id", evt.ID.String()).
		Str("status", string(event.StatusPending)).
		Msg("event created and processed")

	return evt, nil
}

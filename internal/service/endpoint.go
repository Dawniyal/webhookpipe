package service

import (
	"context"

	"github.com/Dawniyal/webhookpipe/internal/middleware"
	"github.com/Dawniyal/webhookpipe/internal/model/endpoint"
	"github.com/Dawniyal/webhookpipe/internal/repository"
)

type EndpointService struct {
	endpointRepo *repository.EndpointsRepository
}

func NewEndpointService(endpointRepo *repository.EndpointsRepository) *EndpointService {
	return &EndpointService{
		endpointRepo: endpointRepo,
	}
}

func (s *EndpointService) AddEndpoint(ctx context.Context, payload *endpoint.AddEndpointPayload) (*endpoint.Endpoint, error) {
	logger := middleware.GetLogger(ctx)

	ep, err := s.endpointRepo.AddEndpoint(ctx, payload)
	if err != nil {
		logger.Error().Err(err).Msg("endpoint validation faild")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "endpoint_created").
		Str("endpoint_id", ep.ID).
		Str("endpoint_url", ep.TargetURL).
		Bool("endpoint_active", ep.Active).
		Msg("endpoint created successfully")

	return ep, nil
}

func (s *EndpointService) GetEndpoint(ctx context.Context, payload *endpoint.GetEndpointByIDPayload) (*endpoint.Endpoint, error) {
	logger := middleware.GetLogger(ctx)

	ep, err := s.endpointRepo.GetEndpointByID(ctx, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch endpoint")
		return nil, err
	}

	return ep, nil
}

func (s *EndpointService) UpdateEndpoint(ctx context.Context, payload *endpoint.UpdateEndpointPayload) (*endpoint.Endpoint, error) {
	logger := middleware.GetLogger(ctx)

	ep, err := s.endpointRepo.UpdateEndpoint(ctx, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to update endpoint")
		return nil, err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "endpoint_updated").
		Str("endpoint_id", ep.ID).
		Str("endpoint_url", ep.TargetURL).
		Bool("endpoint_active", ep.Active).
		Msg("endpoint updated successfully")

	return ep, nil
}

func (s *EndpointService) DeleteEndpointSoft(ctx context.Context, payload *endpoint.DeleteEndpointByIDPayload) error {
	logger := middleware.GetLogger(ctx)

	err := s.endpointRepo.DeleteEndpointSoft(ctx, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to soft delete endpoint")
		return err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "endpoint_deleted_soft").
		Str("endpoint_id", payload.ID).
		Msg("Category soft deleted successfully")

	return nil
}

func (s *EndpointService) DeleteEndpointHard(ctx context.Context, payload *endpoint.DeleteEndpointByIDPayload) error {
	logger := middleware.GetLogger(ctx)

	err := s.endpointRepo.DeleteEndpointHard(ctx, payload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to hard delete endpoint")
		return err
	}

	eventLogger := middleware.GetLogger(ctx)
	eventLogger.Info().
		Str("event", "endpoint_deleted_hard").
		Str("endpoint_id", payload.ID).
		Msg("Category deleted successfully")

	return nil
}

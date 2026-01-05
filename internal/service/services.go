package service

import (
	"github.com/Dawniyal/webhookpipe/internal/lib/job"
	"github.com/Dawniyal/webhookpipe/internal/repository"
)

type Services struct {
	Event    *EventService
	Endpoint *EndpointService
	job      *job.JobService
}

func NewServices(job *job.JobService, repos *repository.Repositories) *Services {
	return &Services{
		Event:    NewEventService(job, repos.Events),
		Endpoint: NewEndpointService(repos.Endpoints),
		job:      job,
	}
}

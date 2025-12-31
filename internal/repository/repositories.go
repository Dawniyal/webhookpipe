package repository

import (
	"github.com/Dawniyal/webhookpipe/internal/server"
)

type Repositories struct {
	Events    *EventsRepository
	Endpoints *EndpointsRepository
}

func NewRepositories(s *server.Server) *Repositories {
	return &Repositories{
		Events:    NewEventsRepository(s),
		Endpoints: NewEndpointsRepository(s),
	}
}

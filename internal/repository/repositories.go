package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	Events    *EventsRepository
	Endpoints *EndpointsRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Events:    NewEventsRepository(db),
		Endpoints: NewEndpointsRepository(db),
	}
}

var ErrEndpointNotFound = errors.New("endpoint not found")

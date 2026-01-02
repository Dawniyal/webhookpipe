package repository

import (
	"context"
	"github.com/Dawniyal/webhookpipe/internal/model/event"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventsRepository struct {
	db *pgxpool.Pool
}

func NewEventsRepository(db *pgxpool.Pool) *EventsRepository {
	return &EventsRepository{db: db}
}

func (r *EventsRepository) CreateEvent(ctx context.Context, e *event.Event) (*event.Event, error) {
	sql := `
		INSERT INTO event (id, endpoint_id, payload, status, created_at)
		VALUES (@id, @endpoint_id, @payload, @status, @created_at)
		RETURNING *;
	`
	args := pgx.NamedArgs{
		"id":          e.ID,
		"endpoint_id": e.EndpointID,
		"payload":     e.Payload,
		"status":      e.Status,
		"created_at":  e.CreatedAt,
	}

	rows, _ := r.db.Query(ctx, sql, args)
	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[event.Event])
	if err != nil {
		return nil, err
	}

	return &res, nil
}

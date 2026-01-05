package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/Dawniyal/webhookpipe/internal/model/event"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventsRepository struct {
	db *pgxpool.Pool
}

func NewEventsRepository(db *pgxpool.Pool) *EventsRepository {
	return &EventsRepository{db: db}
}

func (r *EventsRepository) CreateEvent(ctx context.Context, e *event.CreateEventPayload) (*event.Event, error) {
	ID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	sql := `
		INSERT INTO event (id, endpoint_id, payload)
		VALUES (@id, @endpoint_id, @payload)
		RETURNING *;
	`
	args := pgx.NamedArgs{
		"id":          ID,
		"endpoint_id": e.EndpointID,
		"payload":     e.Payload,
	}

	rows, err := r.db.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[event.Event])
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *EventsRepository) GetEventByIDPayload(ctx context.Context, e *event.GetEventByIDPayload) (*event.Event, error) {
	sql := `SELECT id,endpoint_id,payload,status,active FROM event WHERE id = $1 LIMIT 1;`

	rows, err := r.db.Query(ctx, sql, e.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[event.Event])
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *EventsRepository) UpdateEvent(ctx context.Context, e *event.UpdateEventPayload) (*event.Event, error) {
	sql := `UPDATE event SET `

	args := pgx.NamedArgs{
		"id": e.ID,
	}

	setClauses := []string{}

	if e.EndpointID != nil {
		setClauses = append(setClauses, "endpoint_id = @endpoint_id")
		args["endpoint_id"] = e.EndpointID
	}

	if e.Payload != nil {
		setClauses = append(setClauses, "payload = @payload")
		args["payload"] = e.Payload
	}

	if e.Status != nil {
		setClauses = append(setClauses, "status = @status")
		args["status"] = e.Status
	}

	if e.Active != nil {
		setClauses = append(setClauses, "active = @active")
		args["active"] = e.Active
	}

	if len(setClauses) == 0 {
		return nil, errors.New("bad request")
	}

	sql += strings.Join(setClauses, ", ")

	sql += ` WHERE id = @id RETURNING id, endpoint_id, payload, status, active, created_at`

	rows, _ := r.db.Query(ctx, sql, args)
	ev, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[event.Event])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &ev, nil
}

func (r *EventsRepository) UpdateStatus(ctx context.Context, eventID uuid.UUID, status event.EventStatus) error {
	sql := `UPDATE event SET status = @status WHERE id= @id`
	args := pgx.NamedArgs{
		"id":     eventID,
		"status": status,
	}

	ct, err := r.db.Exec(ctx, sql, args)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("0 Row Affected")
	}

	return nil
}

func (r *EventsRepository) DeleteEventSoft(ctx context.Context, payload *event.DeleteEventPayload) error {
	sql := `UPDATE event SET active = false WHERE id = $1;`
	ct, err := r.db.Exec(ctx, sql, payload.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *EventsRepository) DeleteEventHard(ctx context.Context, payload *event.DeleteEventPayload) error {
	sql := `DELETE FROM event WHERE id = $1;`
	ct, err := r.db.Exec(ctx, sql, payload.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

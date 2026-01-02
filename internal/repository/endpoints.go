package repository

import (
	"context"
	"errors"

	"github.com/Dawniyal/webhookpipe/internal/model/endpoint"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EndpointsRepository struct {
	db *pgxpool.Pool
}

func NewEndpointsRepository(db *pgxpool.Pool) *EndpointsRepository {
	return &EndpointsRepository{db: db}
}

func (r *EndpointsRepository) AddEndpoint(ctx context.Context, payload *endpoint.AddEndpointPayload) (*endpoint.Endpoint, error) {
	sql := `
		INSERT INTO endpoint (id, target_url, active) 
		VALUES (@id, @target_url, @active)
		RETURNING id, target_url, active, created_at;
	`
	args := pgx.NamedArgs{
		"id":         payload.ID,
		"target_url": payload.TargetURL,
		"active":     payload.Active,
	}

	rows, _ := r.db.Query(ctx, sql, args)
	ep, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[endpoint.Endpoint])
	if err != nil {
		return nil, err
	}

	return &ep, nil
}

func (r *EndpointsRepository) GetTargetURLByID(ctx context.Context, payload *endpoint.GetEndpointByIDPayload) (string, error) {
	sql := `SELECT target_url FROM endpoint WHERE id = $1 AND active = TRUE LIMIT 1;`

	var url string
	err := r.db.QueryRow(ctx, sql, payload.ID).Scan(&url)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrEndpointNotFound
		}
		return "", err
	}

	return url, nil
}

func (r *EndpointsRepository) UpdateTargetURLEndpoint(ctx context.Context, payload *endpoint.UpdateEndpointPayload) (*endpoint.Endpoint, error) {
	sql := `
		UPDATE endpoint
		SET target_url = @target_url, active = @active
		WHERE id = @id
		RETURNING id, target_url, active, created_at;
	`
	args := pgx.NamedArgs{
		"id":         payload.ID,
		"target_url": payload.TargetURL,
		"active":     payload.Active,
	}

	rows, _ := r.db.Query(ctx, sql, args)
	ep, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[endpoint.Endpoint])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrEndpointNotFound
		}
		return nil, err
	}

	return &ep, nil
}

func (r *EndpointsRepository) DeleteEndpointSoft(ctx context.Context, payload *endpoint.DeleteEndpointByIDPayload) error {
	sql := `UPDATE endpoint SET active = false WHERE id = $1;`
	ct, err := r.db.Exec(ctx, sql, payload.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrEndpointNotFound
	}
	return nil
}

func (r *EndpointsRepository) DeleteEndpointHard(ctx context.Context, payload *endpoint.DeleteEndpointByIDPayload) error {
	sql := `DELETE FROM endpoint WHERE id = $1;`
	ct, err := r.db.Exec(ctx, sql, payload.ID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return ErrEndpointNotFound
	}
	return nil
}

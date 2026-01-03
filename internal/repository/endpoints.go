package repository

import (
	"context"
	"errors"
	"strings"

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
		RETURNING id, target_url, active, created_at
		LIMIT 1;
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
			return "", pgx.ErrNoRows
		}
		return "", err
	}

	return url, nil
}

func (r *EndpointsRepository) UpdateTargetURLEndpoint(ctx context.Context, payload *endpoint.UpdateEndpointPayload) (*endpoint.Endpoint, error) {
	sql := `UPDATE endpoint SET `

	args := pgx.NamedArgs{
		"id": payload.ID,
	}

	setClauses := []string{}

	if payload.TargetURL != nil {
		setClauses = append(setClauses, "target_url = @target_url")
		args["target_url"] = payload.TargetURL
	}

	if payload.Active != nil {
		setClauses = append(setClauses, "active = @active")
		args["active"] = payload.Active
	}

	if len(setClauses) == 0 {
		return nil, errors.New("bad request") // need to make a err package
	}

	sql += strings.Join(setClauses, ", ")

	sql += ` WHERE id = @id RETURNING id, target_url, active, created_at`

	rows, _ := r.db.Query(ctx, sql, args)
	ep, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[endpoint.Endpoint])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
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
		return pgx.ErrNoRows
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
		return pgx.ErrNoRows
	}
	return nil
}

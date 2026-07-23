package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/model"
)

type UserRepository struct {
	querier db.Querier
}

func NewUserRepository(querier db.Querier) *UserRepository {
	return &UserRepository{querier}
}

func (r *UserRepository) List(ctx context.Context) ([]model.User, error) {
	sql := `SELECT * FROM users WHERE deleted_at IS NULL`
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
}

func (r *UserRepository) Get(ctx context.Context, id model.Id) (*model.User, error) {
	sql := `SELECT * FROM users WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.User])
}

func (r *UserRepository) Add(ctx context.Context, displayName string) (*model.User, error) {
	sql := `INSERT INTO users (display_name) VALUES (@display_name) RETURNING *`
	args := pgx.StrictNamedArgs{"display_name": displayName}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.User])
}

func (r *UserRepository) Delete(ctx context.Context, id model.Id) error {
	sql := `UPDATE users SET deleted_at = CURRENT_TIMESTAMP WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id}
	_, err := r.querier.Exec(ctx, sql, args)
	return err
}

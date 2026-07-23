package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/model"
)

type UserSPIRepository struct {
	querier db.Querier
}

func NewUserSPIRepository(querier db.Querier) *UserSPIRepository {
	return &UserSPIRepository{querier}
}

func (r *UserSPIRepository) List(ctx context.Context) ([]model.UserSPI, error) {
	sql := `SELECT * FROM users_spi WHERE deleted_at IS NULL`
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.UserSPI])
}

func (r *UserSPIRepository) Get(ctx context.Context, id model.Id) (*model.UserSPI, error) {
	sql := `SELECT * FROM users_spi WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.UserSPI])
}

func (r *UserSPIRepository) Delete(ctx context.Context, id model.Id) error {
	sql := `UPDATE users_spi SET deleted_at = CURRENT_TIMESTAMP WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id}
	_, err := r.querier.Exec(ctx, sql, args)
	return err
}

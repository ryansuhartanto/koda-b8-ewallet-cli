package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/model"
)

type TransactionRepository struct {
	querier db.Querier
}

func NewTransactionRepository(querier db.Querier) *TransactionRepository {
	return &TransactionRepository{querier}
}

func (r *TransactionRepository) List(ctx context.Context) ([]model.Transaction, error) {
	sql := `SELECT * FROM transactions WHERE deleted_at IS NULL`
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Transaction])
}

func (r *TransactionRepository) Get(ctx context.Context, id model.Id) (*model.Transaction, error) {
	sql := `SELECT * FROM transactions WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Transaction])
}

func (r *TransactionRepository) Delete(ctx context.Context, id model.Id) error {
	sql := `UPDATE transactions SET deleted_at = CURRENT_TIMESTAMP WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id}
	_, err := r.querier.Exec(ctx, sql, args)
	return err
}

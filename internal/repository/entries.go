package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/model"
)

type EntryRepository struct {
	querier db.Querier
}

func NewEntryRepository(querier db.Querier) *EntryRepository {
	return &EntryRepository{querier}
}

func (r *EntryRepository) List(ctx context.Context) ([]model.Entry, error) {
	sql := `SELECT * FROM entries`
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Entry])
}

func (r *EntryRepository) Get(ctx context.Context, idWallet, idTransaction model.Id) (*model.Entry, error) {
	sql := `SELECT * FROM entries WHERE id_wallet = @id_wallet AND id_transaction = @id_transaction`
	args := pgx.StrictNamedArgs{"id_wallet": idWallet, "id_transaction": idTransaction}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Entry])
}

func (r *EntryRepository) Delete(ctx context.Context, idWallet, idTransaction model.Id) error {
	sql := `DELETE FROM entries WHERE id_wallet = @id_wallet AND id_transaction = @id_transaction`
	args := pgx.StrictNamedArgs{"id_wallet": idWallet, "id_transaction": idTransaction}
	_, err := r.querier.Exec(ctx, sql, args)
	return err
}

func (r *EntryRepository) Add(ctx context.Context, idWallet, idTransaction model.Id, amount, balanceAfter int64) (*model.Entry, error) {
	sql := `
		INSERT INTO entries (id_wallet, id_transaction, amount, balance_idr_after)
		VALUES (@id_wallet, @id_transaction, @amount, @balance_idr_after)
		RETURNING *
	`
	args := pgx.StrictNamedArgs{
		"id_wallet":         idWallet,
		"id_transaction":    idTransaction,
		"amount":            amount,
		"balance_idr_after": balanceAfter,
	}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Entry])
}

func (r *EntryRepository) ListByWallet(ctx context.Context, idWallet model.Id) ([]model.EntryDetail, error) {
	sql := `
		SELECT entries.*, transactions.type, transactions.status, transactions.ref_internal, transactions.note
		FROM entries
		JOIN transactions ON transactions.id = entries.id_transaction
		WHERE entries.id_wallet = @id_wallet
		ORDER BY entries.created_at DESC
	`
	args := pgx.StrictNamedArgs{"id_wallet": idWallet}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.EntryDetail])
}

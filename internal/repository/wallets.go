package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/model"
)

type WalletRepository struct {
	querier db.Querier
}

func NewWalletRepository(querier db.Querier) *WalletRepository {
	return &WalletRepository{querier}
}

func (r *WalletRepository) List(ctx context.Context) ([]model.Wallet, error) {
	sql := `SELECT * FROM wallets WHERE deleted_at IS NULL`
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Wallet])
}

func (r *WalletRepository) Get(ctx context.Context, id model.Id) (*model.Wallet, error) {
	sql := `SELECT * FROM wallets WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Wallet])
}

func (r *WalletRepository) Delete(ctx context.Context, id model.Id) error {
	sql := `UPDATE wallets SET deleted_at = CURRENT_TIMESTAMP WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id}
	_, err := r.querier.Exec(ctx, sql, args)
	return err
}

func (r *WalletRepository) Add(ctx context.Context) (*model.Wallet, error) {
	sql := `INSERT INTO wallets DEFAULT VALUES RETURNING *`
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Wallet])
}

// GetForUpdate locks the wallet row; call within a transaction before mutating its balance.
func (r *WalletRepository) GetForUpdate(ctx context.Context, id model.Id) (*model.Wallet, error) {
	sql := `SELECT * FROM wallets WHERE id = @id AND deleted_at IS NULL FOR UPDATE`
	args := pgx.StrictNamedArgs{"id": id}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.Wallet])
}

func (r *WalletRepository) UpdateBalance(ctx context.Context, id model.Id, balance int64) error {
	sql := `UPDATE wallets SET balance_idr = @balance_idr WHERE id = @id AND deleted_at IS NULL`
	args := pgx.StrictNamedArgs{"id": id, "balance_idr": balance}
	_, err := r.querier.Exec(ctx, sql, args)
	return err
}

func (r *WalletRepository) ListWithOwner(ctx context.Context) ([]model.WalletWithOwner, error) {
	sql := `
		SELECT wallets.*, users.display_name
		FROM wallets
		JOIN users_wallets ON users_wallets.id_wallet = wallets.id
		JOIN users ON users.id = users_wallets.id_user
		WHERE wallets.deleted_at IS NULL
		ORDER BY wallets.id
	`
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.WalletWithOwner])
}

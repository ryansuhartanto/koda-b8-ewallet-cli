package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/model"
)

type UserWalletRepository struct {
	querier db.Querier
}

func NewUserWalletRepository(querier db.Querier) *UserWalletRepository {
	return &UserWalletRepository{querier}
}

func (r *UserWalletRepository) List(ctx context.Context) ([]model.UserWallet, error) {
	sql := `SELECT * FROM users_wallets`
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.UserWallet])
}

func (r *UserWalletRepository) Get(ctx context.Context, idUser, idWallet model.Id) (*model.UserWallet, error) {
	sql := `SELECT * FROM users_wallets WHERE id_user = @id_user AND id_wallet = @id_wallet`
	args := pgx.StrictNamedArgs{"id_user": idUser, "id_wallet": idWallet}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[model.UserWallet])
}

func (r *UserWalletRepository) Delete(ctx context.Context, idUser, idWallet model.Id) error {
	sql := `DELETE FROM users_wallets WHERE id_user = @id_user AND id_wallet = @id_wallet`
	args := pgx.StrictNamedArgs{"id_user": idUser, "id_wallet": idWallet}
	_, err := r.querier.Exec(ctx, sql, args)
	return err
}

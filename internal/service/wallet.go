package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/model"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/repository"
)

var ErrInsufficientBalance = errors.New("service: insufficient balance")
var ErrSameWallet = errors.New("service: source and destination wallet are the same")

// ponytail: timestamp-based, swap for a proper ID generator if concurrent txns collide
func newRef(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, strconv.FormatInt(time.Now().UnixNano(), 36))
}

func noteOrNil(note string) *string {
	if note == "" {
		return nil
	}
	return &note
}

func Register(ctx context.Context, pool *pgxpool.Pool, displayName string) (*model.User, *model.Wallet, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	users := repository.NewUserRepository(tx)
	wallets := repository.NewWalletRepository(tx)
	userWallets := repository.NewUserWalletRepository(tx)

	user, err := users.Add(ctx, displayName)
	if err != nil {
		return nil, nil, err
	}

	wallet, err := wallets.Add(ctx)
	if err != nil {
		return nil, nil, err
	}

	if _, err := userWallets.Add(ctx, user.Id, wallet.Id); err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return user, wallet, nil
}

func ListWallets(ctx context.Context, pool *pgxpool.Pool) ([]model.WalletWithOwner, error) {
	return repository.NewWalletRepository(pool).ListWithOwner(ctx)
}

func History(ctx context.Context, pool *pgxpool.Pool, walletID model.Id) ([]model.EntryDetail, error) {
	return repository.NewEntryRepository(pool).ListByWallet(ctx, walletID)
}

func TopUp(ctx context.Context, pool *pgxpool.Pool, walletID model.Id, amount int64, note string) (*model.Entry, error) {
	return move(ctx, pool, walletID, amount, model.TransactionTypeTopup, newRef("topup"), note)
}

func Withdraw(ctx context.Context, pool *pgxpool.Pool, walletID model.Id, amount int64, note string) (*model.Entry, error) {
	return move(ctx, pool, walletID, -amount, model.TransactionTypeWithdraw, newRef("withdraw"), note)
}

func Payment(ctx context.Context, pool *pgxpool.Pool, walletID model.Id, amount int64, note string) (*model.Entry, error) {
	return move(ctx, pool, walletID, -amount, model.TransactionTypePayment, newRef("payment"), note)
}

// move applies a signed amount to a single wallet as one transaction + entry.
func move(
	ctx context.Context,
	pool *pgxpool.Pool,
	walletID model.Id,
	signedAmount int64,
	typ model.TransactionType,
	ref string,
	note string,
) (*model.Entry, error) {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	wallets := repository.NewWalletRepository(tx)
	transactions := repository.NewTransactionRepository(tx)
	entries := repository.NewEntryRepository(tx)

	wallet, err := wallets.GetForUpdate(ctx, walletID)
	if err != nil {
		return nil, err
	}

	newBalance := wallet.BalanceIdr + signedAmount
	if newBalance < 0 {
		return nil, ErrInsufficientBalance
	}

	if err := wallets.UpdateBalance(ctx, walletID, newBalance); err != nil {
		return nil, err
	}

	transaction, err := transactions.Add(ctx, typ, model.TransactionStatusSuccess, ref, noteOrNil(note))
	if err != nil {
		return nil, err
	}

	entry, err := entries.Add(ctx, walletID, transaction.Id, signedAmount, newBalance)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return entry, nil
}

func Transfer(ctx context.Context, pool *pgxpool.Pool, fromID, toID model.Id, amount int64, note string) (fromEntry, toEntry *model.Entry, err error) {
	if fromID == toID {
		return nil, nil, ErrSameWallet
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback(ctx)

	wallets := repository.NewWalletRepository(tx)
	transactions := repository.NewTransactionRepository(tx)
	entries := repository.NewEntryRepository(tx)

	// Lock both wallets in a fixed (ascending ID) order so two transfers in
	// opposite directions can't deadlock against each other.
	loID, hiID := fromID, toID
	if hiID < loID {
		loID, hiID = hiID, loID
	}
	loWallet, err := wallets.GetForUpdate(ctx, loID)
	if err != nil {
		return nil, nil, err
	}
	hiWallet, err := wallets.GetForUpdate(ctx, hiID)
	if err != nil {
		return nil, nil, err
	}

	from, to := loWallet, hiWallet
	if loID != fromID {
		from, to = hiWallet, loWallet
	}

	if from.BalanceIdr < amount {
		return nil, nil, ErrInsufficientBalance
	}

	fromBalance := from.BalanceIdr - amount
	toBalance := to.BalanceIdr + amount

	if err := wallets.UpdateBalance(ctx, fromID, fromBalance); err != nil {
		return nil, nil, err
	}
	if err := wallets.UpdateBalance(ctx, toID, toBalance); err != nil {
		return nil, nil, err
	}

	transaction, err := transactions.Add(ctx, model.TransactionTypeTransfer, model.TransactionStatusSuccess, newRef("transfer"), noteOrNil(note))
	if err != nil {
		return nil, nil, err
	}

	fromEntry, err = entries.Add(ctx, fromID, transaction.Id, -amount, fromBalance)
	if err != nil {
		return nil, nil, err
	}

	toEntry, err = entries.Add(ctx, toID, transaction.Id, amount, toBalance)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, err
	}

	return fromEntry, toEntry, nil
}

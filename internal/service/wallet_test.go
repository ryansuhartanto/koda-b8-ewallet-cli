package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/service"
)

func TestWalletFlow(t *testing.T) {
	_ = godotenv.Load("../../.env")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "")
	if err != nil {
		t.Skipf("no database available: %v", err)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		t.Skipf("no database available: %v", err)
	}

	m, err := db.Migrate(pool)
	if err != nil {
		t.Fatalf("migrate: %v", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("migrate up: %v", err)
	}

	alice, aliceWallet, err := service.Register(ctx, pool, "Alice Test")
	if err != nil {
		t.Fatalf("register alice: %v", err)
	}
	bob, bobWallet, err := service.Register(ctx, pool, "Bob Test")
	if err != nil {
		t.Fatalf("register bob: %v", err)
	}

	topupEntry, err := service.TopUp(ctx, pool, aliceWallet.Id, 100_000, "initial topup")
	if err != nil {
		t.Fatalf("topup: %v", err)
	}

	fromEntry, toEntry, err := service.Transfer(ctx, pool, aliceWallet.Id, bobWallet.Id, 30_000, "gift")
	if err != nil {
		t.Fatalf("transfer: %v", err)
	}

	t.Cleanup(func() {
		for _, txID := range []any{topupEntry.IDTransaction, fromEntry.IDTransaction} {
			pool.Exec(ctx, `DELETE FROM entries WHERE id_transaction = @id`, pgx.StrictNamedArgs{"id": txID})
			pool.Exec(ctx, `DELETE FROM transactions WHERE id = @id`, pgx.StrictNamedArgs{"id": txID})
		}
		pool.Exec(ctx, `DELETE FROM users_wallets WHERE id_user IN (@a, @b)`, pgx.StrictNamedArgs{"a": alice.Id, "b": bob.Id})
		pool.Exec(ctx, `DELETE FROM wallets WHERE id IN (@w1, @w2)`, pgx.StrictNamedArgs{"w1": aliceWallet.Id, "w2": bobWallet.Id})
		pool.Exec(ctx, `DELETE FROM users WHERE id IN (@a, @b)`, pgx.StrictNamedArgs{"a": alice.Id, "b": bob.Id})
	})

	wallets, err := service.ListWallets(ctx, pool)
	if err != nil {
		t.Fatalf("list wallets: %v", err)
	}

	var aliceBalance, bobBalance int64 = -1, -1
	for _, w := range wallets {
		switch w.Id {
		case aliceWallet.Id:
			aliceBalance = w.BalanceIdr
		case bobWallet.Id:
			bobBalance = w.BalanceIdr
		}
	}

	if aliceBalance != 70_000 {
		t.Fatalf("alice balance = %d, want 70000", aliceBalance)
	}
	if bobBalance != 30_000 {
		t.Fatalf("bob balance = %d, want 30000", bobBalance)
	}

	if _, err := service.Withdraw(ctx, pool, aliceWallet.Id, 1_000_000, ""); !errors.Is(err, service.ErrInsufficientBalance) {
		t.Fatalf("withdraw over balance: got %v, want ErrInsufficientBalance", err)
	}

	if _, _, err := service.Transfer(ctx, pool, aliceWallet.Id, aliceWallet.Id, 1, ""); !errors.Is(err, service.ErrSameWallet) {
		t.Fatalf("transfer to self: got %v, want ErrSameWallet", err)
	}

	history, err := service.History(ctx, pool, aliceWallet.Id)
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(history) != 2 { // topup + transfer-out
		t.Fatalf("history len = %d, want 2", len(history))
	}

	_ = toEntry
}

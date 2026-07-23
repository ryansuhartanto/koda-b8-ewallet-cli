package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"charm.land/huh/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/model"
	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/service"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Panicln("Error loading .env file", err)
	}
}

var ctx context.Context
var pool *pgxpool.Pool

func init() {
	var err error
	ctx = context.Background()
	pool, err = pgxpool.New(ctx, "")
	if err != nil {
		log.Panicln("Error connecting to database", err)
	}

	m, err := db.Migrate(pool)
	if err != nil {
		log.Panicln("Error creating migration", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Panicln("Error migrating database", err)
	}
}

func main() {
	defer pool.Close()

	for {
		action, err := mainMenu()
		if errors.Is(err, huh.ErrUserAborted) || action == "quit" {
			return
		}
		if err != nil {
			log.Fatalln(err)
		}

		if err := dispatch(action); err != nil {
			fmt.Println("Error:", err)
		}
		fmt.Println()
	}
}

func mainMenu() (string, error) {
	var action string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("E-Wallet").
				Options(
					huh.NewOption("Register", "register"),
					huh.NewOption("Top up", "topup"),
					huh.NewOption("Withdraw", "withdraw"),
					huh.NewOption("Transfer between wallets", "transfer"),
					huh.NewOption("Make a payment", "payment"),
					huh.NewOption("View wallets", "wallets"),
					huh.NewOption("View wallet history", "history"),
					huh.NewOption("Quit", "quit"),
				).
				Value(&action),
		),
	).RunWithContext(ctx)
	return action, err
}

func dispatch(action string) error {
	switch action {
	case "register":
		return doRegister()
	case "topup":
		return doMove("Top up", service.TopUp)
	case "withdraw":
		return doMove("Withdraw", service.Withdraw)
	case "payment":
		return doMove("Payment", service.Payment)
	case "transfer":
		return doTransfer()
	case "wallets":
		return doListWallets()
	case "history":
		return doHistory()
	}
	return nil
}

func doRegister() error {
	var displayName string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Display name").
				Validate(nonEmpty).
				Value(&displayName),
		),
	).RunWithContext(ctx)
	if err != nil {
		return err
	}

	user, wallet, err := service.Register(ctx, pool, displayName)
	if err != nil {
		return err
	}

	fmt.Printf("Registered %q as user #%d with wallet #%d\n", user.DisplayName, user.Id, wallet.Id)
	return nil
}

type moveFunc func(ctx context.Context, pool *pgxpool.Pool, walletID model.Id, amount int64, note string) (*model.Entry, error)

func doMove(label string, fn moveFunc) error {
	walletID, ok, err := pickWallet(label + " - pick a wallet")
	if err != nil || !ok {
		return err
	}

	amount, note, err := askAmountAndNote()
	if err != nil {
		return err
	}

	entry, err := fn(ctx, pool, walletID, amount, note)
	if err != nil {
		return err
	}

	fmt.Printf("%s of Rp %d recorded. New balance: Rp %d\n", label, amount, entry.BalanceIdrAfter)
	return nil
}

func doTransfer() error {
	fromID, ok, err := pickWallet("Transfer from")
	if err != nil || !ok {
		return err
	}

	toID, ok, err := pickWallet("Transfer to")
	if err != nil || !ok {
		return err
	}

	amount, note, err := askAmountAndNote()
	if err != nil {
		return err
	}

	fromEntry, toEntry, err := service.Transfer(ctx, pool, fromID, toID, amount, note)
	if err != nil {
		return err
	}

	fmt.Printf("Transferred Rp %d. Wallet #%d balance: Rp %d, wallet #%d balance: Rp %d\n",
		amount, fromID, fromEntry.BalanceIdrAfter, toID, toEntry.BalanceIdrAfter)
	return nil
}

func doListWallets() error {
	wallets, err := service.ListWallets(ctx, pool)
	if err != nil {
		return err
	}
	if len(wallets) == 0 {
		fmt.Println("No wallets yet.")
		return nil
	}

	for _, w := range wallets {
		fmt.Printf("#%d %s - Rp %d\n", w.Id, w.DisplayName, w.BalanceIdr)
	}
	return nil
}

func doHistory() error {
	walletID, ok, err := pickWallet("View history for which wallet?")
	if err != nil || !ok {
		return err
	}

	entries, err := service.History(ctx, pool, walletID)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("No transactions yet.")
		return nil
	}

	for _, e := range entries {
		note := ""
		if e.Note != nil {
			note = ": " + *e.Note
		}
		fmt.Printf("%s  %-8s  %-12s  balance Rp %-12s%s\n",
			e.CreatedAt.Format(time.RFC3339),
			e.TransactionType,
			fmt.Sprintf("%+d", e.Amount),
			fmt.Sprintf("%+d", e.BalanceIdrAfter),
			note,
		)
	}
	return nil
}

func pickWallet(title string) (model.Id, bool, error) {
	wallets, err := service.ListWallets(ctx, pool)
	if err != nil {
		return 0, false, err
	}
	if len(wallets) == 0 {
		fmt.Println("No wallets.")
		return 0, false, nil
	}

	opts := make([]huh.Option[model.Id], len(wallets))
	for i, w := range wallets {
		opts[i] = huh.NewOption(
			fmt.Sprintf("#%d %s - Rp %d", w.Id, w.DisplayName, w.BalanceIdr),
			w.Id,
		)
	}

	var id model.Id
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[model.Id]().Title(title).Options(opts...).Value(&id),
		),
	).RunWithContext(ctx)
	if errors.Is(err, huh.ErrUserAborted) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}

	return id, true, nil
}

func askAmountAndNote() (int64, string, error) {
	var amountStr, note string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Amount (IDR)").
				Validate(validAmount).
				Value(&amountStr),
			huh.NewInput().
				Title("Note (optional)").
				Value(&note),
		),
	).RunWithContext(ctx)
	if err != nil {
		return 0, "", err
	}

	amount, _ := strconv.ParseInt(strings.TrimSpace(amountStr), 10, 64)
	return amount, note, nil
}

func validAmount(s string) error {
	n, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return errors.New("enter a whole number")
	}
	if n <= 0 {
		return errors.New("must be greater than zero")
	}
	return nil
}

func nonEmpty(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("required")
	}
	return nil
}

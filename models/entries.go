package models

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/db"
)

type Entry struct {
	CreatedAt *time.Time

	IDWallet      db.ID
	IDTransaction db.ID

	Amount          int64
	BalanceIdrAfter int64
}

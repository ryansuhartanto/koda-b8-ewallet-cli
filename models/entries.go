package models

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/db"
)

type Entry struct {
	CreatedAt *time.Time

	IdWallet      db.ID
	IdTransaction db.ID

	Amount          int64
	BalanceIdrAfter int64
}

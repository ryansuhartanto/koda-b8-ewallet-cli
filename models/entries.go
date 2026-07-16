package models

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/db"
)

type Entry struct {
	db.ID

	CreatedAt *time.Time

	IdWallet      int64
	IdTransaction int64

	Amount          int64
	BalanceIdrAfter int64
}

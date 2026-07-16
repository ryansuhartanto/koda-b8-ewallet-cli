package models

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/db"
)

type Wallet struct {
	db.ID

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	BalanceIdr int64
}

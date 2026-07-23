package model

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/internal/db"
)

type Wallet struct {
	db.ID

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	BalanceIdr int64
}

type RepoWallets db.RepoSoftDelete[Wallet, db.ID]

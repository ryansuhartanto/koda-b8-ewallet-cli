package models

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/db"
)

type UserWallet struct {
	CreatedAt *time.Time

	IDUser   db.ID
	IDWallet db.ID
}

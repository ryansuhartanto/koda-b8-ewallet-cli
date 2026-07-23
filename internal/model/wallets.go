package model

import "time"

type Wallet struct {
	Id `db:"id"`

	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	BalanceIdr int64 `db:"balance_idr"`
}

type WalletWithOwner struct {
	Wallet

	DisplayName string `db:"display_name"`
}

package model

import "time"

type Entry struct {
	CreatedAt *time.Time `db:"created_at"`

	IDWallet      Id `db:"id_wallet"`
	IDTransaction Id `db:"id_transaction"`

	Amount          int64 `db:"amount"`
	BalanceIdrAfter int64 `db:"balance_idr_after"`
}

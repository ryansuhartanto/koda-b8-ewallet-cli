package models

import "time"

type Entry struct {
	Id int64

	CreatedAt *time.Time

	IdWallet      int64
	IdTransaction int64

	Amount          int64
	BalanceIdrAfter int64
}

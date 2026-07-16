package models

import "time"

type Wallet struct {
	Id int64

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	BalanceIdr int64
}

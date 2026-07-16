package models

import "time"

type UserWallet struct {
	Id int64

	CreatedAt *time.Time
	DeletedAt *time.Time

	IdUser   int64
	IdWallet int64
}

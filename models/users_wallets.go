package models

import "time"

type UserWallet struct {
	CreatedAt *time.Time

	IdUser   int64
	IdWallet int64
}

package model

import "time"

type UserWallet struct {
	CreatedAt *time.Time `db:"created_at"`

	IDUser   Id `db:"id_user"`
	IDWallet Id `db:"id_wallet"`
}

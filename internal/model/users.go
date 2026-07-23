package model

import "time"

type User struct {
	Id `db:"id"`

	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	DisplayName string `db:"display_name"`
}

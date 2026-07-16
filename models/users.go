package models

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/db"
)

type User struct {
	db.ID

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	DisplayName string
}

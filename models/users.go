package models

import "time"

type User struct {
	Id int64

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	DisplayName string
}

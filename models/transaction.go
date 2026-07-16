package models

import "time"

type Transaction struct {
	Id int64

	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time

	Type   TransactionType
	Status TransactionStatus

	RefInternal *string
	RefExternal *string
	Provider    *string
	Note        *string
}

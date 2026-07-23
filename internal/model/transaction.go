package model

import "time"

type TransactionType string

const (
	TransactionTypeTopup    TransactionType = "topup"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypePayment  TransactionType = "payment"
)

type TransactionStatus string

const (
	TransactionStatusPending TransactionStatus = "pending"
	TransactionStatusSuccess TransactionStatus = "success"
	TransactionStatusFailed  TransactionStatus = "failed"
)

type Transaction struct {
	Id `db:"id"`

	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`

	Type   TransactionType   `db:"type"`
	Status TransactionStatus `db:"status"`

	RefInternal *string `db:"ref_internal"`
	RefExternal *string `db:"ref_external"`
	Provider    *string `db:"provider"`
	Note        *string `db:"note"`
}

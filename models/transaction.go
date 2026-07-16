package models

import (
	"time"

	"github.com/ryansuhartanto/koda-b8-ewallet-cli/db"
)

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
	db.ID

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

type RepoTransactions db.RepoSoftDelete[Transaction, db.ID]

package models

type TransactionType string

const (
	TransactionTypeTopup    TransactionType = "topup"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypePayment  TransactionType = "payment"
)

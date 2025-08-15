package domain

import "context"

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction Transaction) error
	GetTransactionByID(ctx context.Context, id string) (*Transaction, error)
	UpdateTransaction(ctx context.Context, transaction Transaction) error
	DeleteTransaction(ctx context.Context, id string) error
}

type TransactionType string

const (
	TransactionTypeCredit     TransactionType = "credit"
	TransactionTypeDebit      TransactionType = "debit"
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeReceipt    TransactionType = "receipt"
	TransactionTypeRefund     TransactionType = "refund"
)

type Transaction struct {
	ID              string          `json:"id"`
	WalletID        string          `json:"wallet_id"`
	Amount          float64         `json:"amount"`
	TransactionType TransactionType `json:"transaction_type"`
	Description     string          `json:"description,omitempty"`
	ReferenceID     string          `json:"reference_id,omitempty"`
	CreatedAt       string          `json:"created_at"`
}

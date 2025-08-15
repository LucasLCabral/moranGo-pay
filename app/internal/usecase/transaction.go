package usecase

import (
	"context"
	"errors"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
)

type TransactionUseCase struct {
	transactionRepo domain.TransactionRepository
}

func NewTransactionUseCase(transactionRepo domain.TransactionRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
	}
}

func (u *TransactionUseCase) CreateTransaction(ctx context.Context, transaction domain.Transaction) error {
	if transaction.WalletID == "" {
		return errors.New("walletID is required")
	}

	if transaction.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if transaction.TransactionType == domain.TransactionType("") {
		return errors.New("transaction type is required")
	}

	return u.transactionRepo.CreateTransaction(ctx, transaction)
}

func (u *TransactionUseCase) GetTransactionByID(ctx context.Context, id string) (*domain.Transaction, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return u.transactionRepo.GetTransactionByID(ctx, id)
}
func (u *TransactionUseCase) UpdateTransaction(ctx context.Context, transaction domain.Transaction) error {
	if transaction.ID == "" {
		return errors.New("id is required")
	}

	return u.transactionRepo.UpdateTransaction(ctx, transaction)
}
func (u *TransactionUseCase) DeleteTransaction(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	return u.transactionRepo.DeleteTransaction(ctx, id)
}

package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/LucasLCabral/moranGo-pay/internal/dto"
	"github.com/google/uuid"
)

type TransactionUseCase struct {
	transactionRepo domain.TransactionRepository
	walletRepo      domain.WalletRepository
}

func NewTransactionUseCase(transactionRepo domain.TransactionRepository, walletRepo domain.WalletRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		walletRepo:      walletRepo,
	}
}

func (u *TransactionUseCase) CreateTransaction(ctx context.Context, transaction domain.Transaction) error {
	if transaction.WalletID == "" {
		return errors.New("walletID is required")
	}

	if transaction.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if transaction.TransactionType == "" {
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

// GetTransactionsByUserID busca todas as transações de um usuário específico
func (u *TransactionUseCase) GetTransactionsByUserID(ctx context.Context, userID string, limit int32) ([]domain.Transaction, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	// Primeiro, buscar a wallet do usuário para obter o walletID
	wallet, err := u.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil || wallet == nil {
		return nil, errors.New("wallet not found")
	}

	// Usar o método do repositório para buscar transações por walletID
	transactions, err := u.transactionRepo.GetTransactionsByWalletID(ctx, wallet.WalletID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	return transactions, nil
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

func (u *TransactionUseCase) calculateAmount(amount float64, txType domain.TransactionType) float64 {
	switch txType {
	case domain.TransactionTypeDeposit, domain.TransactionTypeReceipt, domain.TransactionTypeRefund:
		return amount
	case domain.TransactionTypeWithdrawal, domain.TransactionTypePayment:
		return -amount
	default:
		return amount
	}
}

func (u *TransactionUseCase) ProcessTransaction(ctx context.Context, req dto.TransactionRequest) error {
	if err := u.validateTransactionRequest(req); err != nil {
		return err
	}

	wallet, err := u.walletRepo.GetWalletByUserID(ctx, req.UserID)
	if err != nil || wallet == nil {
		return errors.New("wallet not found")
	}

	if req.Type == domain.TransactionTypeWithdrawal ||
		req.Type == domain.TransactionTypePayment ||
		req.Type == domain.TransactionTypeTransfer {

		if wallet.Balance < req.Amount {
			return errors.New("insufficient balance")
		}
	}

	finalAmount := u.calculateAmount(req.Amount, req.Type)

	wallet.Balance += finalAmount
	wallet.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")

	if err := u.walletRepo.UpdateWallet(ctx, *wallet); err != nil {
		return err
	}

	transaction := &domain.Transaction{
		ID:              uuid.New().String(),
		WalletID:        wallet.WalletID,
		Amount:          finalAmount,
		TransactionType: req.Type,
		Description:     req.Description,
		ReferenceID:     uuid.New().String(),
		CreatedAt:       time.Now().Format("2006-01-02T15:04:05Z"),
	}

	return u.transactionRepo.CreateTransaction(ctx, *transaction)
}

func (u *TransactionUseCase) validateTransactionRequest(req dto.TransactionRequest) error {
	if req.UserID == "" {
		return errors.New("userID is required")
	}

	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	return nil
}

package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/google/uuid"
)

type WalletUseCase struct {
	walletRepo      domain.WalletRepository
	transactionRepo domain.TransactionRepository
}

func NewWalletUseCase(walletRepo domain.WalletRepository, transactionRepo domain.TransactionRepository) *WalletUseCase {
	return &WalletUseCase{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *WalletUseCase) CreateWallet(ctx context.Context, userID string) (*domain.Wallet, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	existing, _ := u.walletRepo.GetWalletByUserID(ctx, userID)
	if existing != nil {
		return existing, nil
	}

	wallet := &domain.Wallet{
		WalletID:  uuid.New().String(),
		UserID:    userID,
		Balance:   0.0,
		CreatedAt: time.Now().Format("2006-01-02"),
		UpdatedAt: time.Now().Format("2006-01-02"),
	}

	if err := u.walletRepo.CreateWallet(ctx, *wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (u *WalletUseCase) GetWallet(ctx context.Context, userID string) (*domain.Wallet, error) {
	if userID == "" {
		return nil, errors.New("userID is required")
	}

	wallet, err := u.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	if wallet == nil {
		return nil, errors.New("wallet not found")
	}

	return wallet, nil
}	

func (u *WalletUseCase) UpdateWallet(ctx context.Context, wallet domain.Wallet) error {
	if wallet.WalletID == "" {
		return errors.New("walletID is required")
	}

	return u.walletRepo.UpdateWallet(ctx, wallet)
}

func (u *WalletUseCase) Deposit(ctx context.Context, userID string, amount float64, description string) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	wallet, err := u.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return err
	}

	wallet.Balance += amount
	wallet.UpdatedAt = time.Now().Format("2006-01-02T15:04:05Z")

	if err := u.walletRepo.UpdateWallet(ctx, *wallet); err != nil {
		return err
	}

	transaction := &domain.Transaction{
		ID: uuid.New().String(),
		WalletID: wallet.WalletID,
		Amount: amount,
		TransactionType: domain.TransactionTypeDeposit,
		Description: description,
		ReferenceID: uuid.New().String(),
		CreatedAt: time.Now().Format("2006-01-02T15:04:05Z"),
	}

	if err := u.transactionRepo.CreateTransaction(ctx, *transaction); err != nil {
		return err
	}

	return nil
}

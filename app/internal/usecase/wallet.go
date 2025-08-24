package usecase

import (
	"context"
	"errors"
	"fmt"
	"math"
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

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	wallet := &domain.Wallet{
		WalletID:  uuid.New().String(),
		UserID:    userID,
		Balance:   0.0,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	if err := u.walletRepo.CreateWallet(ctx, *wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

func (u *WalletUseCase) GetWalletByUserID(ctx context.Context, userID string) (*domain.Wallet, error) {
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

func (u *WalletUseCase) ValidateAndUpdateBalance(ctx context.Context, userID string, amount float64) (*domain.Wallet, error) {
	wallet, err := u.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil || wallet == nil {
		return nil, errors.New("wallet not found")
	}

	if amount < 0 && wallet.Balance < math.Abs(amount) {
		return nil, errors.New("insufficient balance")
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")

	wallet.Balance += amount
	wallet.UpdatedAt = currentTime

	if err := u.walletRepo.UpdateWallet(ctx, *wallet); err != nil {
		return nil, err
	}

	return wallet, nil
}

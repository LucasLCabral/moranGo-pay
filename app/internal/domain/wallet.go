package domain

import "context"

type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet Wallet) error
	GetWalletByUserID(ctx context.Context, userID string) (*Wallet, error)
	UpdateWallet(ctx context.Context, wallet Wallet) error
	DeleteWallet(ctx context.Context, id string) error
}

type Wallet struct {
	WalletID        string `json:"wallet_id"`
	UserID    string `json:"user_id"`
	Balance   float64 `json:"balance"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}


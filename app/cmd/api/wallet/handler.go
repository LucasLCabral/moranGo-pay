package wallet

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/LucasLCabral/moranGo-pay/internal/usecase"
	"github.com/aws/aws-lambda-go/events"
)

type WalletHandler struct {
	walletUseCase *usecase.WalletUseCase
}

func NewWalletHandler(walletUseCase *usecase.WalletUseCase) *WalletHandler {
	return &WalletHandler{
		walletUseCase: walletUseCase,
	}
}


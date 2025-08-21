package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

// GET /wallet/{userID}
func (h *WalletHandler) GetWallet(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := request.PathParameters["userID"]
	if userID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "userID is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	wallet, err := h.walletUseCase.GetWallet(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 404,
				Body:       `{"error": "wallet not found"}`,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}
		
		log.Printf("Error getting wallet: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "internal server error"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	responseBody, err := json.Marshal(wallet)
	if err != nil {
		log.Printf("Error marshaling wallet: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "internal server error"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

// POST /wallet/{userID}/deposit
func (h *WalletHandler) Deposit(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := request.PathParameters["userID"]
	if userID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "userID is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var depositRequest struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	err := json.Unmarshal([]byte(request.Body), &depositRequest)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "invalid request body"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	err = h.walletUseCase.Deposit(ctx, userID, depositRequest.Amount, depositRequest.Description)
	if err != nil {
		log.Printf("Error making deposit: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       `{"message": "deposit successful"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
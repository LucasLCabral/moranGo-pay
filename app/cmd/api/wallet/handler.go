package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/LucasLCabral/moranGo-pay/internal/dto"
	"github.com/LucasLCabral/moranGo-pay/internal/usecase"
	"github.com/aws/aws-lambda-go/events"
)

type WalletHandler struct {
	walletUseCase      *usecase.WalletUseCase
	transactionUseCase *usecase.TransactionUseCase
}

func NewWalletHandler(walletUseCase *usecase.WalletUseCase, transactionUseCase *usecase.TransactionUseCase) *WalletHandler {
	return &WalletHandler{
		walletUseCase:      walletUseCase,
		transactionUseCase: transactionUseCase,
	}
}

// TODO(refactor): add a middleware for validation of UserID, Request Body, etc..

// GET /wallet/{userID}
func (h *WalletHandler) GetWalletByUserID(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := request.PathParameters["userID"]
	if userID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "userID is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	wallet, err := h.walletUseCase.GetWalletByUserID(ctx, userID)
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

	var depositRequest dto.TransactionRequest

	err := json.Unmarshal([]byte(request.Body), &depositRequest)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "invalid request body"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	err = h.transactionUseCase.ProcessTransaction(ctx, dto.TransactionRequest{
		UserID:      userID,
		Amount:      depositRequest.Amount,
		Type:        domain.TransactionTypeDeposit,
		Description: depositRequest.Description,
	})
	if err != nil {
		log.Printf("Error processing transaction: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf(`{"message": "%s processed successfully"}`, domain.TransactionTypeDeposit),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

// POST /wallet/{userID}/withdrawal
func (h *WalletHandler) Withdrawal(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := request.PathParameters["userID"]
	if userID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "userID is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var withdrawalRequest dto.TransactionRequest

	err := json.Unmarshal([]byte(request.Body), &withdrawalRequest)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "invalid request body"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	err = h.transactionUseCase.ProcessTransaction(ctx, dto.TransactionRequest{
		UserID:      userID,
		Amount:      withdrawalRequest.Amount,
		Type:        domain.TransactionTypeWithdrawal,
		Description: withdrawalRequest.Description,
	})
	if err != nil {
		log.Printf("Error processing transaction: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf(`{"message": "%s processed successfully"}`, domain.TransactionTypeWithdrawal),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

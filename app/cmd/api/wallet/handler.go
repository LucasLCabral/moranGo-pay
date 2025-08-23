package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

// TODO(refactor): add a separate middleware for validation of UserID, Request Body, etc..

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

// POST /wallet/{userID}/transaction
func (h *WalletHandler) ProcessTransaction(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := request.PathParameters["userID"]
	if userID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "userID is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	if request.Body == "" {

		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "request body is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	contentType := ""
	for key, value := range request.Headers {
		if strings.ToLower(key) == "content-type" {
			contentType = value
			break
		}
	}

	if contentType != "application/json" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "Content-Type must be application/json, received: '%s'"}`, contentType),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var transactionRequest dto.TransactionRequest
	if err := json.Unmarshal([]byte(request.Body), &transactionRequest); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	transactionRequest.UserID = userID

	if err := h.validateTransactionType(transactionRequest.Type); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	err := h.transactionUseCase.ProcessTransaction(ctx, transactionRequest)
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
		Body:       fmt.Sprintf(`{"message": "%s processed successfully"}`, transactionRequest.Type),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (h *WalletHandler) validateTransactionType(transactionType domain.TransactionType) error {
	validTypes := []domain.TransactionType{
		domain.TransactionTypeDeposit,
		domain.TransactionTypeWithdrawal,
		domain.TransactionTypeTransfer,
		domain.TransactionTypePayment,
		domain.TransactionTypeReceipt,
		domain.TransactionTypeRefund,
		domain.TransactionTypeCredit,
		domain.TransactionTypeDebit,
	}

	for _, validType := range validTypes {
		if transactionType == validType {
			return nil
		}
	}

	return fmt.Errorf("invalid transaction type: %s", transactionType)
}

// GET /wallet/{transactionID}/transaction
func (h *WalletHandler) GetTransactionByID(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	transactionID := request.PathParameters["transactionID"]
	if transactionID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "transactionID is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	transaction, err := h.transactionUseCase.GetTransactionByID(ctx, transactionID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 404,
				Body:       `{"error": "transaction not found"}`,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}

		log.Printf("Error getting transaction: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "internal server error"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	responseBody, err := json.Marshal(transaction)
	if err != nil {
		log.Printf("Error marshaling transaction: %v", err)
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

// GET /wallet/{userID}/transactions
func (h *WalletHandler) GetTransactionsByUserID(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	userID := request.PathParameters["userID"]
	if userID == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 400,
			Body:       `{"error": "userID is required"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Limite padrão de 100 transações
	limit := int32(100)

	// Se houver query parameter para limit, usar o valor fornecido
	if limitStr := request.QueryStringParameters["limit"]; limitStr != "" {
		if parsedLimit, err := strconv.ParseInt(limitStr, 10, 32); err == nil && parsedLimit > 0 {
			limit = int32(parsedLimit)
		}
	}

	transactions, err := h.transactionUseCase.GetTransactionsByUserID(ctx, userID, limit)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return events.APIGatewayV2HTTPResponse{
				StatusCode: 404,
				Body:       `{"error": "wallet not found"}`,
				Headers:    map[string]string{"Content-Type": "application/json"},
			}, nil
		}

		log.Printf("Error getting transactions: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       `{"error": "internal server error"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	responseBody, err := json.Marshal(transactions)
	if err != nil {
		log.Printf("Error marshaling transactions: %v", err)
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

package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/LucasLCabral/moranGo-pay/internal/usecase"

	"github.com/aws/aws-lambda-go/events"
)

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (handler *AuthHandler) Login(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if request.Body == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "{error: Request body is empty}",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	// Unmarshal the request body into the LoginRequest struct
	var loginRequest LoginRequest
	if err := json.Unmarshal([]byte(request.Body), &loginRequest); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "{error: Invalid request body}",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	loginResult, err := handler.authUseCase.Login(ctx, usecase.LoginCredentials{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	})

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "{error: \"Invalid credentials\"}",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	response := &LoginResponse{
		AccessToken:  loginResult.AccessToken,
		RefreshToken: loginResult.RefreshToken,
		TokenType:    loginResult.TokenType,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "{error: \"Failed to marshal response\"}",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
        StatusCode: http.StatusOK,
        Body:       string(responseBody),
        Headers:    map[string]string{"Content-Type": "application/json"},
    }, nil
}

func (handler *AuthHandler) Register(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if request.Body == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "{error: Request body is empty}",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var registerReq struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
		Name     string `json:"name,omitempty"`
	}

	if err := json.Unmarshal([]byte(request.Body), &registerReq); err != nil {
		log.Printf("❌ DEBUG: Erro ao deserializar o corpo da requisição: %v", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "{error: Invalid request body}",
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}	

	log.Printf("🔍 DEBUG: Chamando UseCase.Register...")
    err := handler.authUseCase.Register(ctx, domain.User{
		Email: registerReq.Email,
		Name:  registerReq.Name,
	}, registerReq.Password)
    log.Printf("🔍 DEBUG: UseCase.Register retornou: err=%v", err)
    
    if err != nil {
        log.Printf("❌ DEBUG: Erro retornado: %v", err)
        return events.APIGatewayV2HTTPResponse{
            StatusCode: http.StatusInternalServerError,
            Body:       `{"error": "` + err.Error() + `"}`,
            Headers:    map[string]string{"Content-Type": "application/json"},
        }, nil
    }
    
    log.Printf("✅ DEBUG: Sucesso, retornando OK")
    return events.APIGatewayV2HTTPResponse{
        StatusCode: http.StatusOK,
        Body:       `{"message": "User registered successfully"}`,
        Headers:    map[string]string{"Content-Type": "application/json"},
    }, nil
}	

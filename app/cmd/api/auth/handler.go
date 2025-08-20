package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/LucasLCabral/moranGo-pay/internal/usecase"

	"github.com/aws/aws-lambda-go/events"
)

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
			Body:       `{"error": "Request body is empty"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var loginRequest domain.LoginCredentials
	if err := json.Unmarshal([]byte(request.Body), &loginRequest); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request body"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	loginResult, err := handler.authUseCase.Login(ctx, loginRequest)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       `{"error":"Invalid credentials"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	response := &domain.LoginResponse{
		AccessToken:  loginResult.AccessToken,
		RefreshToken: loginResult.RefreshToken,
		TokenType:    loginResult.TokenType,
		IDToken:      loginResult.IDToken,
	}

	responseBody, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error":"Failed to marshal response"}`,
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
			Body:       `{"error":"Request body is empty"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var registerReq struct {
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
		Name     string `json:"name,omitempty"`
		Username string `json:"username,omitempty"`
	}

	if err := json.Unmarshal([]byte(request.Body), &registerReq); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"Invalid request body"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	err := handler.authUseCase.Register(ctx, domain.User{
		Email: registerReq.Email,
		Name:  registerReq.Name,
		Username: registerReq.Username,
	}, registerReq.Password)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "` + err.Error() + `"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message": "User registered successfully"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (handler *AuthHandler) ConfirmUser(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	if request.Body == "" {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"Request body is empty"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	var confirmUserInput domain.ConfirmUserInput
	if err := json.Unmarshal([]byte(request.Body), &confirmUserInput); err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error":"Invalid request body"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	err := handler.authUseCase.ConfirmUser(ctx, confirmUserInput)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "` + err.Error() + `"}`,
			Headers:    map[string]string{"Content-Type": "application/json"},
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Body:       `{"message": "User confirmed successfully"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

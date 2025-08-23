package dto

import "github.com/LucasLCabral/moranGo-pay/internal/domain"

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         domain.User `json:"user"` // Just passed the user object for helps to get the userID on development test, in prd delete this and use the IDToken
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	IDToken      string      `json:"id_token"`
}

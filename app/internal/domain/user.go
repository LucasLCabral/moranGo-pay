package domain

import (
	"context"
)

// Dependencies Inversion Principle
// High-level modules should not depend on low-level modules. Both should depend on abstractions.
// Abstractions should not depend on details. Details should depend on abstractions.

// UserRepository is an abstraction that defines the methods for interacting with the user data store
type UserRepository interface {
	CreateUser(ctx context.Context, user User, password string) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	ValidateCredentials(ctx context.Context, email string, password string) (bool, error)
	AdminConfirmUser(ctx context.Context, name string) error
	Authenticate(ctx context.Context, email, password string) (*CognitoTokens, error)
}

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CognitoTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	IDToken      string `json:"id_token"`
}

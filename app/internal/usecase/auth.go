package usecase

import (
	"context"
	"errors"
	"time"
)

// Dependencies Inversion Principle
// High-level modules should not depend on low-level modules. Both should depend on abstractions.
// Abstractions should not depend on details. Details should depend on abstractions.

// UserRepository is an abstraction that defines the methods for interacting with the user data store
type UserRepository interface {
	CreateUser(ctx context.Context, user User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// TokenService is an abstraction that defines the methods for generating and validating tokens
type TokenService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (bool, error)
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

type AuthUseCase struct {
	userRepo     UserRepository
	tokenService TokenService
}

func NewAuthUseCase(userRepo UserRepository, tokenService TokenService) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Login implements the Login method of the AuthUseCase interface
func (u *AuthUseCase) Login(ctx context.Context, credentials LoginCredentials) (*LoginResult, error) {
	if err := u.validateLoginCredentials(credentials); err != nil {
		return nil, err
	}

	user, err := u.userRepo.GetUserByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !u.validadePassword(credentials.Password, user) {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := u.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	return &LoginResult{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: "refresh_token_placeholder",
		TokenType:    "Bearer",
	}, nil
}

func (u *AuthUseCase) Register(ctx context.Context, user User, password string) error {
	if err := u.validateUserData(user, password); err != nil {
		return err
	}

	// Check if user already exists
	existingUser, _ := u.userRepo.GetUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("user already exists")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func (u *AuthUseCase) validateLoginCredentials(credentials LoginCredentials) error {
	if credentials.Email == "" || credentials.Password == "" {
		return errors.New("email and password are required")
	}

	return nil
}

func (u *AuthUseCase) validateUserData(user User, password string) error {
	if user.Email == "" || user.Name == "" {
		return errors.New("email and name are required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func (u *AuthUseCase) validadePassword(password string, user *User) bool {
	return true
}

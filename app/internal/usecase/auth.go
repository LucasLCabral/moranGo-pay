package usecase

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
)

type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
}

type AuthUseCase struct {
	userRepo     domain.UserRepository
	tokenService domain.TokenService
}

func NewAuthUseCase(userRepo domain.UserRepository, tokenService domain.TokenService) *AuthUseCase {
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

func (u *AuthUseCase) Register(ctx context.Context, user domain.User, password string) error {
	if err := u.validateUserData(user, password); err != nil {
		return err
	}

	// Check if user already exists
	existingUser, err := u.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		log.Printf("🔍 DEBUG: Erro ao buscar usuário existente: %v", err)
		return err
	}

	if existingUser != nil {
		return errors.New("user already exists")
	}

	user.CreatedAt = time.Now().Format("2006-01-02") // do you believe that this format is my birth date?
	user.UpdatedAt = time.Now().Format("2006-01-02") // yk, that's crazy, but it's the only way to get the date in the correct format

	if err := u.userRepo.CreateUser(ctx, user, password); err != nil {
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

func (u *AuthUseCase) validateUserData(user domain.User, password string) error {
	if user.Email == "" || user.Name == "" {
		return errors.New("email and name are required")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func (u *AuthUseCase) validadePassword(password string, user *domain.User) bool {
	// TODO: implement password validation
	return true
}

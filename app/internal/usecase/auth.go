package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
)

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

func (u *AuthUseCase) Register(ctx context.Context, user domain.User, password string) error {
	if err := u.validateUserData(user, password); err != nil {
		return err
	}

	existingUser, err := u.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
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

func (u *AuthUseCase) ConfirmUser(ctx context.Context, in domain.ConfirmUserInput) error {
	if in.Username == "" && in.Email == "" {
		return errors.New("username or email is required")
	}

	username := in.Username
	if username == "" {
		user, err := u.userRepo.GetUserByEmail(ctx, in.Email)
		if err != nil {
			return err
		}
		if user == nil || user.ID == "" {
			return errors.New("user not found")
		}
		username = user.ID
	}

	return u.userRepo.AdminConfirmUser(ctx, username)
}

func (u *AuthUseCase) Login(ctx context.Context, c domain.LoginCredentials) (*domain.LoginResult, error) {
	if err := u.validateLoginCredentials(c); err != nil { return nil, err }

	tok, err := u.userRepo.Authenticate(ctx, c.Email, c.Password)
	if err != nil { return nil, errors.New("invalid credentials") }

	return &domain.LoginResult{
		AccessToken:  tok.AccessToken,
		RefreshToken: tok.RefreshToken,
		IDToken:      tok.IDToken,
		TokenType:    tok.TokenType,
	}, nil
}

func (u *AuthUseCase) validateLoginCredentials(credentials domain.LoginCredentials) error {
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

package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/google/uuid"
)

type AuthUseCase struct {
	userRepo    domain.UserRepository
	profileRepo domain.ProfileRepository
	walletRepo  domain.WalletRepository
}

func NewAuthUseCase(userRepo domain.UserRepository, profileRepo domain.ProfileRepository, walletRepo domain.WalletRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
		walletRepo:  walletRepo,
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

	existingProfile, _ := u.profileRepo.GetProfileByUsername(ctx, user.Username)
	if existingProfile != nil {
		return errors.New("username already taken")
	}

	userID := uuid.New().String()
	user.CreatedAt = time.Now().Format("2006-01-02") // do you believe that this format is my birth date? :)
	user.UpdatedAt = time.Now().Format("2006-01-02") // yk, that's crazy, but it's the only way to get the date in the correct format

	if err := u.userRepo.CreateUser(ctx, user, password); err != nil {
		return errors.New("failed to create user")
	}

	user = domain.User{
		ID:        userID,
		Email:     user.Email,
		Name:      user.Name,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	profile := domain.UserProfile{
		PK:        fmt.Sprintf("USER#%s", userID),
		SK:        "PROFILE",
		GSI1PK:    fmt.Sprintf("EMAIL#%s", user.Email),
		GSI1SK:    fmt.Sprintf("USER#%s", userID),
		UserID:    userID,
		Email:     user.Email,
		Name:      user.Name,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}

	if err := u.profileRepo.CreateProfile(ctx, profile); err != nil {
		// TODO: Implementar rollback do Cognito em caso de falha
		return errors.New("failed to create user profile")
	}

	// Removido usernameIndex redundante - username já está no perfil principal

	wallet := domain.Wallet{
		WalletID:  uuid.New().String(),
		UserID:    userID,
		Balance:   0.0,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err := u.walletRepo.CreateWallet(ctx, wallet); err != nil {
		return errors.New("failed to create wallet")
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
	if err := u.validateLoginCredentials(c); err != nil {
		return nil, err
	}

	tok, err := u.userRepo.Authenticate(ctx, c.Email, c.Password)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	userProfile, err := u.profileRepo.GetUserByEmail(ctx, c.Email)
	if err != nil || userProfile == nil {
		return nil, errors.New("user profile not found")
	}

	return &domain.LoginResult{
		User: domain.User{
			ID:       userProfile.UserID,
			Email:    userProfile.Email,
			Name:     userProfile.Name,
			Username: userProfile.Username,
		},
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
	if user.Email == "" || user.Name == "" || user.Username == "" {
		return errors.New("email, name and username are required")
	}

	if len(user.Username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}

	if len(user.Username) > 20 {
		return errors.New("username must be at most 20 characters long")
	}

	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

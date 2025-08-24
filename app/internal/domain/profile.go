package domain

import "context"

type ProfileRepository interface {
	CreateProfile(ctx context.Context, profile UserProfile) error
	GetProfileByUserID(ctx context.Context, userID string) (*UserProfile, error)
	GetProfileByUsername(ctx context.Context, username string) (*UserProfile, error)
	GetUserByEmail(ctx context.Context, email string) (*UserProfile, error)

	// TODO: Implement these methods
	UpdateProfile(ctx context.Context, profile UserProfile) error
	DeleteProfile(ctx context.Context, userID string) error
}

type UserProfile struct {
	PK        string `json:"pk"`
	SK        string `json:"sk"`
	Type      string `json:"type"`
	GSI1PK    string `json:"gsi1pk"`
	GSI1SK    string `json:"gsi1sk"`
	GSI2PK    string `json:"gsi2pk"`
	GSI2SK    string `json:"gsi2sk"`
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

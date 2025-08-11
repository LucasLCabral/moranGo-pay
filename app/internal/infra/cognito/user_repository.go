package cognito

import (
	"context"
	"fmt"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type CognitoUserRepository struct {
	cognitoClient *cognitoidentityprovider.Client
	userPoolID    string
	appClientID   string
}

func NewCognitoUserRepository(userPoolID, appClientID string) (domain.UserRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Criar cliente Cognito
	cognitoClient := cognitoidentityprovider.NewFromConfig(cfg)

	return &CognitoUserRepository{
		userPoolID:    userPoolID,
		appClientID:   appClientID,
		cognitoClient: cognitoClient,
	}, nil
}

func (r *CognitoUserRepository) CreateUser(ctx context.Context, user domain.User, password string) error {
	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(r.appClientID),
		Username: aws.String(user.Email),
		Password: aws.String(password),
		UserAttributes: []types.AttributeType{
			{
				Name: aws.String("email"), Value: aws.String(user.Email),
			},
			{
				Name: aws.String("name"), Value: aws.String(user.Name),
			},
		},
	}

	_, err := r.cognitoClient.SignUp(ctx, signUpInput)
	if err != nil {
		return fmt.Errorf("failed to sign up user: %w", err)
	}

	return nil
}

func (r *CognitoUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	input := &cognitoidentityprovider.ListUsersInput{
		UserPoolId: aws.String(r.userPoolID),
		Filter:     aws.String(fmt.Sprintf("email = \"%s\"", email)),
		Limit:      aws.Int32(1),
	}

	result, err := r.cognitoClient.ListUsers(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	if len(result.Users) == 0 {
		return nil, nil
	}

	u := result.Users[0]

	var nameAttr, emailAttr string
	for _, a := range u.Attributes {
		if aws.ToString(a.Name) == "email" {
			emailAttr = aws.ToString(a.Value)
		} else if aws.ToString(a.Name) == "name" {
			nameAttr = aws.ToString(a.Value)
		}
	}

	return &domain.User{
		ID:        aws.ToString(u.Username),
		Email:     emailAttr,
		Name:      nameAttr,
		CreatedAt: u.UserCreateDate.Format("2006-01-02"),
		UpdatedAt: u.UserLastModifiedDate.Format("2006-01-02"),
	}, nil
}

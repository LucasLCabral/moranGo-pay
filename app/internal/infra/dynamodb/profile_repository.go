// app/internal/infra/dynamodb/profile_repository.go
package dynamodb

import (
	"context"
	"fmt"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type profileItem struct {
	PK        string `dynamodbav:"PK"`
	SK        string `dynamodbav:"SK"`
	GSI1PK    string `dynamodbav:"GSI1PK"`
	GSI1SK    string `dynamodbav:"GSI1SK"`
	UserID    string `dynamodbav:"user_id"`
	Email     string `dynamodbav:"email"`
	Name      string `dynamodbav:"name"`
	Username  string `dynamodbav:"username"`
	CreatedAt string `dynamodbav:"created_at"`
}

type DynamoProfileRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoProfileRepository(tableName string) (domain.ProfileRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoProfileRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (r *DynamoProfileRepository) CreateProfile(ctx context.Context, profile domain.UserProfile) error {
	item := profileItem{
		PK:        profile.PK,
		SK:        profile.SK,
		GSI1PK:    fmt.Sprintf("USERNAME#%s", profile.Username),
		GSI1SK:    fmt.Sprintf("USER#%s", profile.UserID),
		UserID:    profile.UserID,
		Email:     profile.Email,
		Name:      profile.Name,
		Username:  profile.Username,
		CreatedAt: profile.CreatedAt,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("failed to put profile: %w", err)
	}

	return nil
}

func (r *DynamoProfileRepository) GetProfileByUserID(ctx context.Context, userID string) (*domain.UserProfile, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			"SK": &types.AttributeValueMemberS{Value: "PROFILE"},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	if result.Item == nil {
		return nil, nil 
	}

	var item profileItem
	err = attributevalue.UnmarshalMap(result.Item, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return &domain.UserProfile{
		PK:        item.PK,
		SK:        item.SK,
		UserID:    item.UserID,
		Email:     item.Email,
		Name:      item.Name,
		Username:  item.Username,
		CreatedAt: item.CreatedAt,
	}, nil
}

func (r *DynamoProfileRepository) GetProfileByUsername(ctx context.Context, username string) (*domain.UserProfile, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("GSI1"),
		KeyConditionExpression: aws.String("GSI1PK = :username"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: fmt.Sprintf("USERNAME#%s", username)},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query profile by username: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil 
	}

	var item profileItem
	err = attributevalue.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return &domain.UserProfile{
		PK:        item.PK,
		SK:        item.SK,
		UserID:    item.UserID,
		Email:     item.Email,
		Name:      item.Name,
		Username:  item.Username,
		CreatedAt: item.CreatedAt,
	}, nil
}

func (r *DynamoProfileRepository) UpdateProfile(ctx context.Context, profile domain.UserProfile) error {
	// TODO: Implementar update
	return fmt.Errorf("UpdateProfile not implemented yet")
}

func (r *DynamoProfileRepository) DeleteProfile(ctx context.Context, userID string) error {
	// TODO: Implementar delete
	return fmt.Errorf("DeleteProfile not implemented yet")
}
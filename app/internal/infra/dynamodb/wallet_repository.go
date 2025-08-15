package dynamodb

import (
	"context"
	"fmt"
	"time"

	"github.com/LucasLCabral/moranGo-pay/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Single Table Design Keys:
// Wallet: PK="USER#<userID>", SK="WALLET#<walletID>"
// GSI1: GSI1PK="USER#<userID>", GSI1SK="WALLET#<walletID>"

type walletItem struct {
	PK        string  `dynamodbav:"PK"`
	SK        string  `dynamodbav:"SK"`
	WalletID  string  `dynamodbav:"wallet_id"`
	UserID    string  `dynamodbav:"user_id"`
	Balance   float64 `dynamodbav:"balance"`
	CreatedAt string  `dynamodbav:"created_at"`
	UpdatedAt string  `dynamodbav:"updated_at"`
}

type DynamoWalletRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoWalletRepository(tableName string) (domain.WalletRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoWalletRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (r *DynamoWalletRepository) CreateWallet(ctx context.Context, wallet domain.Wallet) error {
	item := walletItem{
		PK:        fmt.Sprintf("USER#%s", wallet.UserID),
		SK:        fmt.Sprintf("WALLET#%s", wallet.WalletID),
		WalletID:  wallet.WalletID,
		UserID:    wallet.UserID,
		Balance:   wallet.Balance,
		CreatedAt: wallet.CreatedAt,
		UpdatedAt: wallet.UpdatedAt,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal wallet: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put wallet: %w", err)
	}

	return nil
}

func (r *DynamoWalletRepository) GetWalletByUserID(ctx context.Context, userID string) (*domain.Wallet, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			":sk": &types.AttributeValueMemberS{Value: "WALLET#"},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query wallet: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	var item walletItem
	err = attributevalue.UnmarshalMap(result.Items[0], &item)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallet: %w", err)
	}

	return &domain.Wallet{
		WalletID:  item.WalletID,
		UserID:    item.UserID,
		Balance:   item.Balance,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (r *DynamoWalletRepository) UpdateWallet(ctx context.Context, wallet domain.Wallet) error {
	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", wallet.UserID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("WALLET#%s", wallet.WalletID)},
		},
		UpdateExpression: aws.String("SET balance = :balance, updated_at = :updated_at"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":balance":    &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", wallet.Balance)},
			":updated_at": &types.AttributeValueMemberS{Value: time.Now().Format("2006-01-02T15:04:05Z")},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	return nil
}

func (r *DynamoWalletRepository) DeleteWallet(ctx context.Context, wallet domain.Wallet) error {
	_, err := r.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", wallet.UserID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("WALLET#%s", wallet.WalletID)},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}

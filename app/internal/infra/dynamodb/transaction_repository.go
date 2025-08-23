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
// Transaction: PK="WALLET#<walletID>", SK="TRANSACTION#<timestamp>#<transactionID>"
// GSI1: GSI1PK="USER#<userID>", GSI1SK="TRANSACTION#<timestamp>#<transactionID>"
// GSI2: GSI2PK="TRANSACTION#<transactionID>", GSI2SK="<timestamp>"

type transactionItem struct {
	PK              string  `dynamodbav:"PK"`
	SK              string  `dynamodbav:"SK"`
	Type            string  `dynamodbav:"type"`
	GSI1PK          string  `dynamodbav:"GSI1PK"`
	GSI1SK          string  `dynamodbav:"GSI1SK"`
	GSI2PK          string  `dynamodbav:"GSI2PK"`
	GSI2SK          string  `dynamodbav:"GSI2SK"`
	ID              string  `dynamodbav:"id"`
	WalletID        string  `dynamodbav:"wallet_id"`
	Amount          float64 `dynamodbav:"amount"`
	TransactionType string  `dynamodbav:"transaction_type"`
	Description     string  `dynamodbav:"description"`
	ReferenceID     string  `dynamodbav:"reference_id"`
	CreatedAt       string  `dynamodbav:"created_at"`
}

type DynamoTransactionRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoTransactionRepository(tableName string) (domain.TransactionRepository, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return &DynamoTransactionRepository{
		client:    client,
		tableName: tableName,
	}, nil
}

func (r *DynamoTransactionRepository) CreateTransaction(ctx context.Context, transaction domain.Transaction) error {
	timestamp := time.Now().Format("2006-01-02T15:04:05Z")

	item := transactionItem{
		PK:              fmt.Sprintf("WALLET#%s", transaction.WalletID),
		SK:              fmt.Sprintf("TRANSACTION#%s#%s", timestamp, transaction.ID),
		Type:            "TRANSACTION",
		GSI1PK:          fmt.Sprintf("USER#%s", ""), // TODO: resolver UserID
		GSI1SK:          fmt.Sprintf("TRANSACTION#%s#%s", timestamp, transaction.ID),
		GSI2PK:          fmt.Sprintf("TRANSACTION#%s", transaction.ID),
		GSI2SK:          timestamp,
		ID:              transaction.ID,
		WalletID:        transaction.WalletID,
		Amount:          transaction.Amount,
		TransactionType: string(transaction.TransactionType),
		Description:     transaction.Description,
		ReferenceID:     transaction.ReferenceID,
		CreatedAt:       transaction.CreatedAt,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      av,
	})

	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}

	return nil
}

func (r *DynamoTransactionRepository) GetTransactionByID(ctx context.Context, id string) (*domain.Transaction, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String("GSI2"),
		KeyConditionExpression: aws.String("GSI2PK = :gsi2pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":gsi2pk": &types.AttributeValueMemberS{Value: fmt.Sprintf("TRANSACTION#%s", id)},
		},
		Limit: aws.Int32(1),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query transaction: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, fmt.Errorf("transaction not found")
	}

	var txItem transactionItem
	err = attributevalue.UnmarshalMap(result.Items[0], &txItem)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}

	transaction := &domain.Transaction{
		ID:              txItem.ID,
		WalletID:        txItem.WalletID,
		Amount:          txItem.Amount,
		TransactionType: domain.TransactionType(txItem.TransactionType),
		Description:     txItem.Description,
		ReferenceID:     txItem.ReferenceID,
		CreatedAt:       txItem.CreatedAt,
	}

	return transaction, nil
}

func (r *DynamoTransactionRepository) UpdateTransaction(ctx context.Context, transaction domain.Transaction) error {
	// TODO: implementar update
	return fmt.Errorf("update transaction not implemented yet")
}

func (r *DynamoTransactionRepository) DeleteTransaction(ctx context.Context, id string) error {
	// TODO: implementar delete
	return fmt.Errorf("delete transaction not implemented yet")
}

// Método adicional útil
func (r *DynamoTransactionRepository) GetTransactionsByWalletID(ctx context.Context, walletID string, limit int32) ([]domain.Transaction, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk_prefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: fmt.Sprintf("WALLET#%s", walletID)},
			":sk_prefix": &types.AttributeValueMemberS{Value: "TRANSACTION#"},
		},
		Limit:            aws.Int32(limit),
		ScanIndexForward: aws.Bool(false), // Mais recentes primeiro
	})

	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}

	var transactions []domain.Transaction
	for _, item := range result.Items {
		var txItem transactionItem
		err = attributevalue.UnmarshalMap(item, &txItem)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal transaction: %w", err)
		}

		transactions = append(transactions, domain.Transaction{
			ID:              txItem.ID,
			WalletID:        txItem.WalletID,
			Amount:          txItem.Amount,
			TransactionType: domain.TransactionType(txItem.TransactionType),
			Description:     txItem.Description,
			ReferenceID:     txItem.ReferenceID,
			CreatedAt:       txItem.CreatedAt,
		})
	}

	return transactions, nil
}

package main

import (
	"context"
	"log"
	"os"

	"github.com/LucasLCabral/moranGo-pay/cmd/api/auth"
	"github.com/LucasLCabral/moranGo-pay/internal/delivery"
	"github.com/LucasLCabral/moranGo-pay/internal/infra/cognito"
	"github.com/LucasLCabral/moranGo-pay/internal/infra/dynamodb"
	"github.com/LucasLCabral/moranGo-pay/internal/usecase"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var router *delivery.Router

func init() {
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	clientID := os.Getenv("COGNITO_CLIENT_ID")
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	if userPoolID == "" || clientID == "" || tableName == "" {
		log.Fatal("Required environment variables not configured!")
	}

	userRepo, err := cognito.NewCognitoUserRepository(userPoolID, clientID)
	if err != nil {
		log.Fatal("Error creating Cognito repository:", err)
	}

	profileRepo, err := dynamodb.NewDynamoProfileRepository(tableName)
	if err != nil {
		log.Fatal("Error creating DynamoDB repository:", err)
	}

	authUseCase := usecase.NewAuthUseCase(userRepo, profileRepo)
	authHandler := auth.NewAuthHandler(authUseCase)

	router = delivery.NewRouter(authHandler)

	log.Println("System initialized successfully! 🚀 🍓🍓🍓")
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: %v", r)
		}
	}()

	response, err := router.Route(ctx, request)
	return response, err
}

func main() {
	lambda.Start(handler)
}

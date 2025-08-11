package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/LucasLCabral/moranGo-pay/cmd/api/auth"
	"github.com/LucasLCabral/moranGo-pay/internal/delivery"
	"github.com/LucasLCabral/moranGo-pay/internal/infra/cognito"
	"github.com/LucasLCabral/moranGo-pay/internal/infra/jwt"
	"github.com/LucasLCabral/moranGo-pay/internal/usecase"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

var router *delivery.Router

func getSSMParameter(paramName string) string {
	if paramName == "" {
		log.Fatal("❌ SSM parameter name is empty!")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	client := ssm.NewFromConfig(cfg)

	resp, err := client.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           &paramName,
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		log.Fatalf("failed to get parameter from SSM: %v", err)
	}

	return *resp.Parameter.Value
}

func init() {
	// 1. Configurar variáveis de ambiente
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID")
	clientID := os.Getenv("COGNITO_CLIENT_ID")
	secretParam := os.Getenv("JWT_SECRET_KEY")

	if userPoolID == "" || clientID == "" || secretParam == "" {
		log.Fatal("❌ Required environment variables not configured!")
	}

	// Buscar JWT Secret seguro no SSM
	secretKey := getSSMParameter(secretParam)

	// 2. Configurar dependências
	userRepo, err := cognito.NewCognitoUserRepository(userPoolID, clientID)
	if err != nil {
		log.Fatal("Error creating Cognito repository:", err)
	}

	tokenService := jwt.NewJWTTokenService(secretKey)
	authUseCase := usecase.NewAuthUseCase(userRepo, tokenService)
	authHandler := auth.NewAuthHandler(authUseCase)

	// 3. Configurar router
	router = delivery.NewRouter(authHandler)

	log.Println("System initialized successfully! 🚀")
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	defer func() {
        if r := recover(); r != nil {
            log.Printf("❌ PANIC capturado: %v", r)
        }
    }()
    
    log.Printf("✅ DEBUG: Handler iniciado")
    response, err := router.Route(ctx, request)
    log.Printf("✅ DEBUG: Handler finalizado")
    return response, err
}

func main() {
	lambda.Start(handler)
}

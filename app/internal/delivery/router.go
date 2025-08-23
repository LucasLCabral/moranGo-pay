package delivery

import (
	"context"
	"strings"

	"github.com/LucasLCabral/moranGo-pay/cmd/api/auth"
	"github.com/LucasLCabral/moranGo-pay/cmd/api/wallet"
	"github.com/aws/aws-lambda-go/events"
)

type Router struct {
	authHandler   *auth.AuthHandler
	walletHandler *wallet.WalletHandler
}

func NewRouter(authHandler *auth.AuthHandler, walletHandler *wallet.WalletHandler) *Router {
	return &Router{
		authHandler:   authHandler,
		walletHandler: walletHandler,
	}
}

func (r *Router) Route(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	path := request.RawPath
	method := request.RequestContext.HTTP.Method

	switch {
	// Auth Routes
	case path == "/auth/login" && method == "POST":
		return r.authHandler.Login(ctx, request)
	case path == "/auth/register" && method == "POST":
		return r.authHandler.Register(ctx, request)
	case path == "/auth/confirm" && method == "POST":
		return r.authHandler.ConfirmUser(ctx, request)

	// Wallet Route (protected by JWT)
	case strings.HasPrefix(path, "/wallet/") && method == "GET":
		// Verificar se é busca por transação específica ou lista de transações
		if strings.HasSuffix(path, "/transaction") {
			return r.walletHandler.GetTransactionByID(ctx, request)
		} else if strings.HasSuffix(path, "/transactions") {
			return r.walletHandler.GetTransactionsByUserID(ctx, request)
		} else {
			// Busca da wallet por userID
			return r.walletHandler.GetWalletByUserID(ctx, request)
		}
	case strings.HasPrefix(path, "/wallet/") && strings.HasSuffix(path, "/transaction") && method == "POST":
		return r.walletHandler.ProcessTransaction(ctx, request)

	default:
		return r.handleNotFound()
	}

}

func (r *Router) handleNotFound() (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 404,
		Body:       `{"error": "Route not found"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

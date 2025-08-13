package delivery

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/LucasLCabral/moranGo-pay/cmd/api/auth"
)

type Router struct {
	authHandler *auth.AuthHandler
}

func NewRouter(authHandler *auth.AuthHandler) *Router {
	return &Router{
		authHandler: authHandler,
	}
}

func (r *Router) Route(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    path := request.RawPath
    method := request.RequestContext.HTTP.Method

    if request.RawPath == "/auth/register" && request.RequestContext.HTTP.Method == "POST" {
        response, err := r.authHandler.Register(ctx, request)
        if err != nil {
            return events.APIGatewayV2HTTPResponse{
                StatusCode: http.StatusInternalServerError,
                Body:       `{"error": "Internal server error"}`,
                Headers:    map[string]string{"Content-Type": "application/json"},
            }, nil
        }
        return response, nil
    }

    switch {
    // Auth Routes
    case path == "/auth/login" && method == "POST":
        return r.authHandler.Login(ctx, request)
    case path == "/auth/register" && method == "POST":
        return r.authHandler.Register(ctx, request)
    case path == "/auth/confirm" && method == "POST":
        return r.authHandler.ConfirmUser(ctx, request)
    // Wallet Route (protected by JWT)
    case path == "/wallet" && method == "GET":
        return r.handleWallet(ctx, request)
    
    default:
        return r.handleNotFound()
    }

}

func (r *Router) handleWallet(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// TODO: Implement wallet handler
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 501, // Not Implemented
		Body:       `{"error": "Wallet endpoint not implemented yet"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func (r *Router) handleNotFound() (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 404,
		Body:       `{"error": "Route not found"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
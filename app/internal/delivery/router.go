package delivery

import (
	"context"
	"log"
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

// ✅ MUDANÇA: Usar APIGatewayV2HTTPRequest e APIGatewayV2HTTPResponse
func (r *Router) Route(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
    // ✅ DEBUG - API Gateway v2 (HTTP API)
    log.Printf("�� DEBUG: RawPath: '%s'", request.RawPath)
    log.Printf("�� DEBUG: HTTP Method: '%s'", request.RequestContext.HTTP.Method)

    // ✅ USAR OS CAMPOS CORRETOS PARA API GATEWAY v2
    path := request.RawPath
    method := request.RequestContext.HTTP.Method

    log.Printf("�� DEBUG: Path final: '%s'", path)
    log.Printf("🔍 DEBUG: Method final: '%s'", method)

    if request.RawPath == "/auth/register" && request.RequestContext.HTTP.Method == "POST" {
        log.Printf("🚀 DEBUG: Chamando AuthHandler.Register...")
        response, err := r.authHandler.Register(ctx, request)
        if err != nil {
            log.Printf("❌ ERRO no AuthHandler.Register: %v", err)
            return events.APIGatewayV2HTTPResponse{
                StatusCode: http.StatusInternalServerError,
                Body:       `{"error": "Internal server error"}`,
                Headers:    map[string]string{"Content-Type": "application/json"},
            }, nil
        }
        log.Printf("✅ AuthHandler.Register executado com sucesso")
        return response, nil
    }

    switch {
    // Auth Routes
    case path == "/auth/login" && method == "POST":
        log.Printf("✅ Rota LOGIN encontrada")
        return r.authHandler.Login(ctx, request)
    case path == "/auth/register" && method == "POST":
        log.Printf("✅ Rota REGISTER encontrada")
        return r.authHandler.Register(ctx, request)
    
    // Wallet Route (protected by JWT)
    case path == "/wallet" && method == "GET":
        log.Printf("✅ Rota WALLET encontrada")
        return r.handleWallet(ctx, request)
    
    default:
        log.Printf("❌ Rota NÃO encontrada: %s %s", method, path)
        return r.handleNotFound()
    }

}

// ✅ MUDANÇA: Atualizar para v2
func (r *Router) handleWallet(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// TODO: Implement wallet handler
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 501, // Not Implemented
		Body:       `{"error": "Wallet endpoint not implemented yet"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

// ✅ MUDANÇA: Atualizar para v2
func (r *Router) handleNotFound() (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 404,
		Body:       `{"error": "Route not found"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
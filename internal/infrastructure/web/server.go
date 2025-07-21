package web

import (
	"ledger-service/internal/core/services/ledger"
	"ledger-service/internal/infrastructure/web/handler/balance"
	"ledger-service/internal/infrastructure/web/handler/transaction"
	"ledger-service/internal/infrastructure/web/middleware"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

type Server struct {
	api                huma.API
	balanceHandler     *balance.Handler
	transactionHandler *transaction.Handler
}

func NewServer(ledgerService *ledger.Service) *Server {
	mux := http.NewServeMux()

	config := huma.DefaultConfig("Ledger API", "1.0.0")
	api := humago.New(mux, config)

	server := &Server{
		api:                api,
		balanceHandler:     balance.NewHandler(ledgerService),
		transactionHandler: transaction.NewHandler(ledgerService),
	}

	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	huma.Register(s.api, huma.Operation{
		OperationID: "create-deposit",
		Method:      http.MethodPost,
		Path:        "/api/customers/{customerId}/transactions/deposits",
		Summary:     "Create a deposit",
		Description: "Create a deposit transaction for a customer. Called by other services when customer adds money.",
		Tags:        []string{"transactions"},
		Errors:      []int{400, 500},
	}, s.transactionHandler.CreateDeposit)

	huma.Register(s.api, huma.Operation{
		OperationID: "create-purchase",
		Method:      http.MethodPost,
		Path:        "/api/customers/{customerId}/transactions/purchase",
		Summary:     "Create a purchase",
		Description: "Create a purchase transaction for a customer. Called by other services when customer buys from restaurant.",
		Tags:        []string{"transactions"},
		Errors:      []int{400, 500},
	}, s.transactionHandler.CreatePurchase)

	huma.Register(s.api, huma.Operation{
		OperationID: "get-balance",
		Method:      http.MethodGet,
		Path:        "/api/balances/{userId}",
		Summary:     "Get user balance",
		Description: "Retrieve the current balance for a specific user. Returns 0 balance for new users.",
		Tags:        []string{"balances"},
		Errors:      []int{500},
	}, s.balanceHandler.GetBalance)

	huma.Register(s.api, huma.Operation{
		OperationID: "get-customer-transactions",
		Method:      http.MethodGet,
		Path:        "/api/customers/{customerId}/transactions",
		Summary:     "Get customer transactions",
		Description: "Retrieve all transactions for a specific customer",
		Tags:        []string{"transactions"},
		Errors:      []int{500},
	}, s.transactionHandler.GetCustomerTransactions)

	huma.Register(s.api, huma.Operation{
		OperationID: "get-restaurant-transactions",
		Method:      http.MethodGet,
		Path:        "/api/restaurants/{restaurantId}/transactions",
		Summary:     "Get restaurant transactions",
		Description: "Retrieve all transactions for a specific restaurant",
		Tags:        []string{"transactions"},
		Errors:      []int{500},
	}, s.transactionHandler.GetRestaurantTransactions)
}

func (s *Server) Handler() http.Handler {
	return middleware.LoggingMiddleware(s.api.Adapter())
}

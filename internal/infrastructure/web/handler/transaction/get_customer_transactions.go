package transaction

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"ledger-service/internal/core/types"
)

type GetCustomerTransactionsInput struct {
	CustomerId string `path:"customerId" doc:"Customer ID"`
}

type GetCustomerTransactionsOutput struct {
	Body []GetCustomerTransactionsResponse `json:"body"`
}

type GetCustomerTransactionsResponse struct {
	Id         string        `json:"id" doc:"Transaction ID"`
	Type       string        `json:"type" doc:"Transaction type"`
	Amount     float32       `json:"amount" doc:"Transaction amount"`
	User       *UserResponse `json:"user,omitempty" doc:"User"`
	Restaurant *UserResponse `json:"restaurant,omitempty" doc:"Restaurant involved in the transaction"`
	CreatedAt  time.Time     `json:"createdAt" doc:"Transaction creation timestamp"`
}

func ToGetCustomerTransactionsResponse(t types.Transaction) GetCustomerTransactionsResponse {
	resp := GetCustomerTransactionsResponse{
		Id:        t.Id,
		Type:      string(t.Type),
		Amount:    t.Amount,
		CreatedAt: t.CreatedAt,
	}

	if t.Customer.Id != "" {
		resp.User = &UserResponse{
			Id:   t.Customer.Id,
			Type: string(t.Customer.Type),
		}
	}

	if t.Restaurant.Id != "" {
		resp.Restaurant = &UserResponse{
			Id:   t.Restaurant.Id,
			Type: string(t.Restaurant.Type),
		}
	}

	return resp
}

func (h *Handler) GetCustomerTransactions(ctx context.Context, input *GetCustomerTransactionsInput) (*GetCustomerTransactionsOutput, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	transactions, err := h.ledgerService.GetCustomerTransactions(ctxWithTimeout, input.CustomerId)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve transactions", err)
	}

	responses := []GetCustomerTransactionsResponse{}
	for _, transaction := range transactions {
		responses = append(responses, ToGetCustomerTransactionsResponse(transaction))
	}

	return &GetCustomerTransactionsOutput{
		Body: responses,
	}, nil
}
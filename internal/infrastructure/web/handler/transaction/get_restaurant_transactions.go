package transaction

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"ledger-service/internal/core/types"
)

type GetRestaurantTransactionsInput struct {
	RestaurantId string `path:"restaurantId" doc:"Restaurant ID"`
}

type GetRestaurantTransactionsOutput struct {
	Body []GetRestaurantTransactionsResponse `json:"body"`
}

type GetRestaurantTransactionsResponse struct {
	Id                 string        `json:"id" doc:"Transaction ID"`
	Type               string        `json:"type" doc:"Transaction type"`
	Amount             float32       `json:"amount" doc:"Transaction amount"`
	Customer           *UserResponse `json:"customer,omitempty" doc:"Customer involved in the transaction"`
	Restaurant         *UserResponse `json:"restaurant,omitempty" doc:"Restaurant involved in the transaction"`
	RelatedTransaction string        `json:"relatedTransaction,omitempty" doc:"Related transaction ID (for commission transactions)"`
	CreatedAt          time.Time     `json:"createdAt" doc:"Transaction creation timestamp"`
}

func ToGetRestaurantTransactionsResponse(t types.Transaction) GetRestaurantTransactionsResponse {
	resp := GetRestaurantTransactionsResponse{
		Id:                 t.Id,
		Type:               string(t.Type),
		Amount:             t.Amount,
		RelatedTransaction: t.RelatedTransaction,
		CreatedAt:          t.CreatedAt,
	}

	if t.Customer.Id != "" {
		resp.Customer = &UserResponse{
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

func (h *Handler) GetRestaurantTransactions(ctx context.Context, input *GetRestaurantTransactionsInput) (*GetRestaurantTransactionsOutput, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	transactions, err := h.ledgerService.GetRestaurantTransactions(ctxWithTimeout, input.RestaurantId)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to retrieve transactions", err)
	}

	responses := []GetRestaurantTransactionsResponse{}
	for _, transaction := range transactions {
		responses = append(responses, ToGetRestaurantTransactionsResponse(transaction))
	}

	return &GetRestaurantTransactionsOutput{
		Body: responses,
	}, nil
}
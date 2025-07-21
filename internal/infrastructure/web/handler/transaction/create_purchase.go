package transaction

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"ledger-service/internal/core/types"
)

type PurchaseRequest struct {
	Amount       float32 `json:"amount" doc:"Purchase amount"`
	RestaurantId string  `json:"restaurantId" doc:"Restaurant ID"`
}

type PurchaseInput struct {
	CustomerId string          `path:"customerId" doc:"Customer ID"`
	Body       PurchaseRequest `json:"body"`
}

type PurchaseOutput struct {
	Body PurchaseResponse `json:"body"`
}

type PurchaseResponse struct {
	Id         string        `json:"id" doc:"Transaction ID"`
	Type       string        `json:"type" doc:"Transaction type (always PURCHASE)"`
	Amount     float32       `json:"amount" doc:"Purchase amount"`
	Customer   *UserResponse `json:"customer" doc:"Customer who made the purchase"`
	Restaurant *UserResponse `json:"restaurant" doc:"Restaurant involved in the purchase"`
	CreatedAt  time.Time     `json:"createdAt" doc:"Transaction creation timestamp"`
}

func (req PurchaseRequest) ToTransaction(customerId string) types.Transaction {
	return types.Transaction{
		Type:   types.PURCHASE,
		Amount: req.Amount,
		Customer: types.User{
			Id:   customerId,
			Type: types.CUSTOMER,
		},
		Restaurant: types.User{
			Id:   req.RestaurantId,
			Type: types.RESTAURANT,
		},
		CreatedAt: time.Now(),
	}
}

func ToPurchaseResponse(t types.Transaction) PurchaseResponse {
	resp := PurchaseResponse{
		Id:        t.Id,
		Type:      string(t.Type),
		Amount:    t.Amount,
		CreatedAt: t.CreatedAt,
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

func (h *Handler) CreatePurchase(ctx context.Context, input *PurchaseInput) (*PurchaseOutput, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	transaction := input.Body.ToTransaction(input.CustomerId)

	id, err := h.ledgerService.SaveTransaction(ctxWithTimeout, transaction)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create purchase", err)
	}

	transaction.Id = id
	response := ToPurchaseResponse(transaction)

	return &PurchaseOutput{
		Body: response,
	}, nil
}
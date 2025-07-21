package transaction

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"ledger-service/internal/core/types"
)

type DepositRequest struct {
	Amount float32 `json:"amount" doc:"Deposit amount"`
}

type DepositInput struct {
	CustomerId string         `path:"customerId" doc:"Customer ID"`
	Body       DepositRequest `json:"body"`
}

type DepositOutput struct {
	Body DepositResponse `json:"body"`
}

type DepositResponse struct {
	Id        string        `json:"id" doc:"Transaction ID"`
	Type      string        `json:"type" doc:"Transaction type (always DEPOSIT)"`
	Amount    float32       `json:"amount" doc:"Deposit amount"`
	Customer  *UserResponse `json:"customer" doc:"Customer who made the deposit"`
	CreatedAt time.Time     `json:"createdAt" doc:"Transaction creation timestamp"`
}

func (req DepositRequest) ToTransaction(customerId string) types.Transaction {
	return types.Transaction{
		Type:   types.DEPOSIT,
		Amount: req.Amount,
		Customer: types.User{
			Id:   customerId,
			Type: types.CUSTOMER,
		},
		CreatedAt: time.Now(),
	}
}

func ToDepositResponse(t types.Transaction) DepositResponse {
	resp := DepositResponse{
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

	return resp
}

func (h *Handler) CreateDeposit(ctx context.Context, input *DepositInput) (*DepositOutput, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	transaction := input.Body.ToTransaction(input.CustomerId)

	id, err := h.ledgerService.SaveTransaction(ctxWithTimeout, transaction)
	if err != nil {
		return nil, huma.Error400BadRequest("Failed to create deposit", err)
	}

	transaction.Id = id
	response := ToDepositResponse(transaction)

	return &DepositOutput{
		Body: response,
	}, nil
}
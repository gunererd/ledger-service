package balance

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"ledger-service/internal/core/types"
)

type GetBalanceInput struct {
	UserId string `path:"userId" doc:"User ID"`
}

type GetBalanceOutput struct {
	Body GetBalanceResponse `json:"body"`
}

type GetBalanceResponse struct {
	UserId          string   `json:"userId" doc:"User ID"`
	Amount          float32  `json:"amount" doc:"Current balance amount"`
	TotalCommission *float32 `json:"totalCommission,omitempty" doc:"Total commission earned (restaurants only)"`
}

func ToGetBalanceResponse(balance types.Balance) GetBalanceResponse {
	response := GetBalanceResponse{
		UserId: balance.UserId,
		Amount: balance.Amount,
	}

	// Include commission if user is a restaurant (has commission > 0)
	if balance.TotalCommission > 0 {
		response.TotalCommission = &balance.TotalCommission
	}

	return response
}

func (h *Handler) GetBalance(ctx context.Context, input *GetBalanceInput) (*GetBalanceOutput, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	balance, err := h.ledgerService.GetBalance(ctxWithTimeout, input.UserId)
	if err != nil {
		return nil, huma.Error500InternalServerError("Internal server error", err)
	}

	response := ToGetBalanceResponse(balance)

	return &GetBalanceOutput{
		Body: response,
	}, nil
}
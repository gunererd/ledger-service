package interfaces

import (
	"context"
	"ledger-service/internal/core/types"
)

type TransactionRepository interface {
	Save(ctx context.Context, t types.Transaction) (string, error)
	GetManyForCustomer(ctx context.Context, id string) ([]types.Transaction, error)
	GetManyForRestaurant(ctx context.Context, id string) ([]types.Transaction, error)
}

type BalanceRepository interface {
	GetBalance(ctx context.Context, userId string) (types.Balance, error)
	UpdateBalance(ctx context.Context, userId string, amount float32) error
	UpdateTotalCommission(ctx context.Context, userId string, amount float32) error
}

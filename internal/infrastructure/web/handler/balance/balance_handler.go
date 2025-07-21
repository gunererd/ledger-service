package balance

import (
	"ledger-service/internal/core/services/ledger"
)

type Handler struct {
	ledgerService *ledger.Service
}

func NewHandler(ledgerService *ledger.Service) *Handler {
	return &Handler{
		ledgerService: ledgerService,
	}
}


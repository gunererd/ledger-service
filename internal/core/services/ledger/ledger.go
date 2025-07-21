package ledger

import (
	"context"
	"ledger-service/internal/core/interfaces"
	"ledger-service/internal/core/types"
	"log/slog"
	"os"
	"time"
)

const COMMISSION_RATE = 0.05

type Service struct {
	transactionRepo interfaces.TransactionRepository
	balanceRepo     interfaces.BalanceRepository
	queue           interfaces.Queue
	ctx             context.Context
	cancel          context.CancelFunc
	logger          *slog.Logger
}

func NewService(transactionRepo interfaces.TransactionRepository, balanceRepo interfaces.BalanceRepository, queue interfaces.Queue) *Service {
	ctx, cancel := context.WithCancel(context.Background())
	service := &Service{
		transactionRepo: transactionRepo,
		balanceRepo:     balanceRepo,
		queue:           queue,
		ctx:             ctx,
		cancel:          cancel,
		logger:          slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	}

	go service.processBalanceUpdates()

	return service
}

func (s *Service) Shutdown() {
	s.cancel()
}

func (s *Service) SaveTransaction(ctx context.Context, transaction types.Transaction) (string, error) {
	id, err := s.transactionRepo.Save(ctx, transaction)
	if err != nil {
		return "", err
	}

	transaction.Id = id
	s.queue.Enqueue(transaction)

	return id, nil
}

func (s *Service) GetBalance(ctx context.Context, userId string) (types.Balance, error) {
	return s.balanceRepo.GetBalance(ctx, userId)
}

func (s *Service) GetCustomerTransactions(ctx context.Context, customerId string) ([]types.Transaction, error) {
	return s.transactionRepo.GetManyForCustomer(ctx, customerId)
}

func (s *Service) GetRestaurantTransactions(ctx context.Context, restaurantId string) ([]types.Transaction, error) {
	return s.transactionRepo.GetManyForRestaurant(ctx, restaurantId)
}

func (s *Service) processBalanceUpdates() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			tx := s.queue.Dequeue()
			if tx.Id == "" {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			if err := s.updateBalances(ctx, tx); err != nil {
				s.logger.Error("Balance update failed",
					"error", err.Error(),
					"transaction_id", tx.Id,
					"transaction_type", string(tx.Type),
					"amount", tx.Amount,
				)
				cancel()
				continue
			}

			s.processCommission(ctx, tx)

			cancel()
		}
	}
}

func (s *Service) updateBalances(ctx context.Context, transaction types.Transaction) error {
	switch transaction.Type {
	case types.DEPOSIT:
		return s.balanceRepo.UpdateBalance(ctx, transaction.Customer.Id, transaction.Amount)
	case types.PURCHASE:
		if err := s.balanceRepo.UpdateBalance(ctx, transaction.Customer.Id, -transaction.Amount); err != nil {
			return err
		}
		return s.balanceRepo.UpdateBalance(ctx, transaction.Restaurant.Id, transaction.Amount)
	case types.COMMISSION:
		// Deduct from current balance
		if err := s.balanceRepo.UpdateBalance(ctx, transaction.Restaurant.Id, -transaction.Amount); err != nil {
			return err
		}
		// Track cumulative commission earned
		return s.balanceRepo.UpdateTotalCommission(ctx, transaction.Restaurant.Id, transaction.Amount)
	}
	return nil
}

func (s *Service) processCommission(ctx context.Context, tx types.Transaction) {
	if !s.shouldApplyCommission(tx) {
		return
	}

	commissionTx := s.buildCommissionTransaction(tx)
	if commissionTx.Amount <= 0 {
		return
	}

	id, err := s.transactionRepo.Save(ctx, commissionTx)
	if err != nil {
		s.logger.Error("Commission transaction save failed",
			"error", err.Error(),
			"original_transaction_id", tx.Id,
			"commission_amount", commissionTx.Amount,
		)
		return
	}

	commissionTx.Id = id
	s.queue.Enqueue(commissionTx)
}

func (s *Service) shouldApplyCommission(tx types.Transaction) bool {
	return tx.Type == types.PURCHASE
}

func (s *Service) buildCommissionTransaction(tx types.Transaction) types.Transaction {
	return types.Transaction{
		Type:               types.COMMISSION,
		Amount:             tx.Amount * COMMISSION_RATE,
		Restaurant:         tx.Restaurant,
		RelatedTransaction: tx.Id,
		CreatedAt:          time.Now(),
	}
}

package loyalty

import (
	"context"
	"time"

	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/internal/orders"
	"github.com/niksmo/gophermart/pkg/logger"
)

var flushInterval = 4 * time.Second

type LoyaltyService struct {
	repository            LoyaltyRepository
	ordersToLoyaltyStream <-chan orders.OrderScheme
}

func NewService(
	ctx context.Context,
	repository LoyaltyRepository,
	ordersToLoyaltyStream <-chan orders.OrderScheme,
) LoyaltyService {
	service := LoyaltyService{
		repository:            repository,
		ordersToLoyaltyStream: ordersToLoyaltyStream,
	}
	go service.flushTransactions(ctx)
	return service
}

func (s LoyaltyService) GetUserBalance(
	ctx context.Context, userID int32,
) (BalanceScheme, error) {
	return s.repository.ReadBalance(ctx, userID)
}

func (s LoyaltyService) WithdrawPoints(
	ctx context.Context, userID int32, orderNumber string, amount float64,
) error {
	return s.repository.ReduceBalance(ctx, userID, orderNumber, amount)
}

func (s LoyaltyService) GetUserWithdrawals(
	ctx context.Context, userID int32,
) (WithdrawalsScheme, error) {
	withdrawals, err := s.repository.ReadWithdrawals(ctx, userID)
	if err != nil {
		return withdrawals, err
	}
	if len(withdrawals) == 0 {
		return withdrawals, errs.ErrLoyaltyNoWithdrawals
	}

	return withdrawals, nil
}

func (s LoyaltyService) flushTransactions(ctx context.Context) {
	log := logger.Instance.With().Caller().Logger()
	ticker := time.NewTicker(flushInterval)
	var transactions []TransactionScheme

	for {
		select {
		case <-ctx.Done():
			return
		case order := <-s.ordersToLoyaltyStream:
			transactions = append(transactions, TransactionScheme{
				UserID:      order.OwnerID,
				OrderNumber: order.Number,
				Amount:      order.Accrual,
			})

			log.Info().
				Str("orderNum", order.Number).
				Float64("amount", order.Accrual).
				Msg("append transaction")
		case <-ticker.C:
			logger.Instance.Info().Caller().Msg("flush transactions tick occur")
			err := s.repository.CreateAddTransactions(ctx, transactions)
			if err != nil {
				log.Error().Err(err).Msg("didn't flush")
				continue
			}
			transactions = nil
		}
	}
}

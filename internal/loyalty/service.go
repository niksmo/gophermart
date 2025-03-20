package loyalty

import (
	"context"

	"github.com/niksmo/gophermart/internal/errs"
)

type LoyaltyService struct {
	repository LoyaltyRepository
}

func NewService(repository LoyaltyRepository) LoyaltyService {
	return LoyaltyService{
		repository: repository,
	}
}

func (s LoyaltyService) GetUserBalance(
	ctx context.Context, userID int32,
) (BalanceScheme, error) {
	return s.repository.ReadBalance(ctx, userID)
}

func (s LoyaltyService) WithdrawPoints(
	ctx context.Context, userID int32, orderNumber string, amount float32,
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

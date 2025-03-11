package loyalty

import "context"

type LoyaltyService struct {
	repository LoyaltyRepository
}

func NewService(repository LoyaltyRepository) LoyaltyService {
	return LoyaltyService{repository: repository}
}

func (s LoyaltyService) GetUsersBalance(
	ctx context.Context, userID int32,
) (BalanceScheme, error) {
	return s.repository.ReadBalance(ctx, userID)
}

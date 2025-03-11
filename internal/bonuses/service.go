package bonuses

import "context"

type BonusesService struct {
	repository BonusesRepository
}

func NewBonusesService(repository BonusesRepository) BonusesService {
	return BonusesService{repository: repository}
}

func (s BonusesService) GetUsersBalance(ctx context.Context, userID int32) (BalanceScheme, error) {
	return s.repository.Read(ctx, userID)
}

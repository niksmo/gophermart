package orders

import (
	"context"
	"errors"

	"github.com/niksmo/gophermart/internal/errs"
)

type OrdersService struct {
	repository OrdersRepository
}

func NewService(repository OrdersRepository) OrdersService {
	return OrdersService{repository: repository}
}

func (s OrdersService) UploadOrder(ctx context.Context, userID int64, orderNumber int64) error {
	err := s.repository.Create(ctx, userID, orderNumber)
	if err != nil {
		if errors.Is(err, errs.ErrOrderUploaded) {
			s.repository.ReadByOrderNumber(ctx, orderNumber)
		}
		return err
	}
	return nil
}

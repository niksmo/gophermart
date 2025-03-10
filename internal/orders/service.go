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

func (s OrdersService) UploadOrder(ctx context.Context, userID int32, orderNumber int64) error {
	err := s.repository.Create(ctx, userID, orderNumber)
	if err != nil {
		if errors.Is(err, errs.ErrOrderUploaded) {
			return s.processUploadConflict(ctx, userID, orderNumber)
		}
		return err
	}
	return nil
}

func (s OrdersService) GetUserOrders(
	ctx context.Context, userID int32, orders []OrderScheme,
) ([]OrderScheme, error) {
	orders, err := s.repository.ReadListByUser(ctx, userID, orders)
	if err != nil {
		return orders, err
	}
	if len(orders) == 0 {
		return orders, errs.ErrOrdersNoUploads
	}

	return orders, nil
}

func (s OrdersService) processUploadConflict(
	ctx context.Context, userID int32, orderNumber int64,
) error {
	order, err := s.repository.ReadByOrderNumber(ctx, orderNumber)
	if err != nil {
		return err
	}
	if order.OwnerID == userID {
		return errs.ErrOrderUploadedByUser
	}
	return errs.ErrOrderUploadedByOther

}

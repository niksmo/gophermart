package orders

import (
	"context"
	"errors"
	"runtime"
	"time"

	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/logger"
)

const (
	REGISTERED = "REGISTERED"
	PROCESSING = "PROCESSING"
	PROCESSED  = "PROCESSED"
	INVALID    = "INVALID"

	pullStreamSize = 1024
)

var flushInterval = 4 * time.Second

type OrdersService struct {
	repository          OrdersRepository
	accrualFetchStream  chan OrderScheme
	accrualResultStream chan AccrualResult
}

func NewService(
	ctx context.Context,
	repository OrdersRepository,
	ordersToLoyaltyStream chan<- OrderScheme,
) OrdersService {
	accrualFetchStream := make(chan OrderScheme, pullStreamSize)
	accrualResultStream := make(chan AccrualResult)

	workerPool := AccrualWorkerPool{
		Num:     runtime.NumCPU(),
		ChanIN:  accrualFetchStream,
		ChanOUT: accrualResultStream,
	}
	workerPool.Run(ctx)

	service := OrdersService{
		repository:          repository,
		accrualFetchStream:  accrualFetchStream,
		accrualResultStream: accrualResultStream,
	}

	go service.flushAccrualResults(ctx, ordersToLoyaltyStream)

	// TO DO restore func

	return service
}

func (s OrdersService) UploadOrder(ctx context.Context, userID int32, orderNumber string) error {
	order, err := s.repository.Create(ctx, userID, orderNumber)
	if err != nil {
		if errors.Is(err, errs.ErrOrderUploaded) {
			return s.handleConflict(ctx, userID, orderNumber)
		}
		return err
	}

	select {
	case s.accrualFetchStream <- order:
	default:
		logger.Instance.Error().
			Str("orderNumber", order.Number).
			Caller().
			Msg("accrual stream is full")
	}
	return nil
}

func (s OrdersService) GetUserOrders(
	ctx context.Context, userID int32,
) (OrderListScheme, error) {
	orderList, err := s.repository.ReadListByUser(ctx, userID)
	if err != nil {
		return orderList, err
	}
	if len(orderList) == 0 {
		return orderList, errs.ErrOrdersNoUploads
	}

	return orderList, nil
}

func (s OrdersService) handleConflict(
	ctx context.Context, userID int32, orderNumber string,
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

func (s OrdersService) flushAccrualResults(
	ctx context.Context, ordersToLoyaltyStream chan<- OrderScheme,
) {
	log := logger.Instance.With().Caller().Logger()
	ticker := time.NewTicker(flushInterval)
	var updatedOrders []OrderScheme

	for {
		select {
		case <-ctx.Done():
			return
		case result := <-s.accrualResultStream:
			if result.Error != nil {
				log.Error().
					Err(result.Error).
					Msg("receive error from result stream")
				continue
			}

			updatedOrders = append(updatedOrders, result.Order)
			log.Info().
				Str("orderNum", result.Order.Number).
				Msg("append order to flush buffer")
		case <-ticker.C:
			log.Info().Msg("flush updated orders tick occur")
			err := s.repository.UpdateAccrual(ctx, updatedOrders)
			if err != nil {
				log.Error().Err(err).Msg("didn't flush")
				continue
			}

			for _, order := range updatedOrders {
				switch order.Status {
				case REGISTERED, PROCESSING:
					log.Info().
						Str("orderNum", order.Number).
						Msg("send to pull stream")

					s.accrualFetchStream <- order
				case PROCESSED:
					log.Info().
						Str("orderNum", order.Number).
						Msg("send to loyalty stream")

					ordersToLoyaltyStream <- order
				}
			}

			updatedOrders = nil
		}
	}
}

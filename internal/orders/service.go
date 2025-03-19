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

	pullStreamSize       = 1024
	maxRestoreListSize   = 100
	restoreBatchInterval = time.Minute
)

var flushInterval = 4 * time.Second

type OrdersService struct {
	repository          OrdersRepository
	accrualFetchStream  chan OrderScheme
	accrualResultStream chan AccrualResult
}

func NewService(
	ctx context.Context, repository OrdersRepository,
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

	return service
}

func (s OrdersService) UploadOrder(
	ctx context.Context, userID int32, orderNumber string,
) error {
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
	return s.determineUploadedErr(userID, order.OwnerID)
}

func (s OrdersService) determineUploadedErr(userID, ownerID int32) error {
	if ownerID == userID {
		return errs.ErrOrderUploadedByUser
	}
	return errs.ErrOrderUploadedByOther
}

func (s OrdersService) FlushAccrualResults(ctx context.Context) {
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
					Msg("receive error from fetch stream")
				continue
			}
			updatedOrders = append(updatedOrders, result.Order)
		case <-ticker.C:
			err := s.repository.UpdateAccrual(ctx, updatedOrders)
			if err != nil {
				log.Error().Err(err).Msg("didn't flush")
				continue
			}

			for _, order := range updatedOrders {
				if order.Status == REGISTERED || order.Status == PROCESSING {
					s.accrualFetchStream <- order
				}
			}
			updatedOrders = nil
		}
	}
}

func (s OrdersService) Restore(ctx context.Context) {
	log := logger.Instance.With().Caller().Logger()
	orders, err := s.repository.ReadNonAccrualList(ctx)
	if err != nil {
		log.Error().Err(err).Msg("didn't restore")
		return
	}

	if len(orders) == 0 {
		return
	}

	log.Info().Int("ordersToRestore", len(orders)).Msg("prepare to restore")

	ticker := time.NewTicker(restoreBatchInterval)
	for {
		size := min(maxRestoreListSize, len(orders))
		for _, order := range orders[:size] {
			select {
			case <-ctx.Done():
				return
			case s.accrualFetchStream <- order:
			}
		}

		orders = orders[size:]
		if len(orders) == 0 {
			break
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}

	log.Info().Msg("all non accrual orders send to fetch stream")
}

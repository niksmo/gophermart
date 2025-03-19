package orders

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/logger"
)

var retryIntervals = [3]time.Duration{
	time.Second,
	3 * time.Second,
	5 * time.Second,
}

type AccrualServiceResponse struct {
	statusCode int
	retryAfter time.Duration
	body       []byte
	err        error
}

func (res AccrualServiceResponse) HasRetryAfter() bool {
	return res.retryAfter == 0
}

func (res AccrualServiceResponse) HasError() bool {
	return res.statusCode == fiber.StatusInternalServerError ||
		res.err != nil
}

func (res AccrualServiceResponse) Err() error {
	if res.err != nil {
		return res.err
	}
	if res.statusCode >= 400 {
		return errs.ErrOrdersAccrualServiceAPI
	}
	return nil
}

type AccrualWorkerPool struct {
	Num     int
	ChanIN  <-chan OrderScheme
	ChanOUT chan<- AccrualResult
}

func (wp AccrualWorkerPool) Run(ctx context.Context, config config.AccrualConfig) {
	for range wp.Num {
		worker := newAccrualWorker(config)
		go worker.Run(ctx, wp.ChanIN, wp.ChanOUT)
	}
}

type AccrualResult struct {
	Error error
	Order OrderScheme
}

type AccrualWorker struct {
	retryIntervals  []time.Duration
	currentInterval time.Duration
	makeRequestURL  func(orderNumber string) (URL string)
}

func newAccrualWorker(config config.AccrualConfig) AccrualWorker {
	return AccrualWorker{
		retryIntervals:  retryIntervals[:],
		makeRequestURL:  config.GetOrdersReqURL,
		currentInterval: retryIntervals[0],
	}
}

func (w *AccrualWorker) Run(
	ctx context.Context,
	orderStream <-chan OrderScheme,
	resultStream chan<- AccrualResult,
) {
	log := logger.Instance.With().Caller().Logger()

	for order := range orderStream {
		select {
		case <-ctx.Done():
			return
		default:
			result := AccrualResult{Order: order}
			data, err := w.getAccrualStatus(ctx, order.Number)
			if err != nil {
				log.Error().
					Err(err).
					Str("orderNum", order.Number).
					Msg("accrual result error")
				result.Error = err
				resultStream <- result
				continue
			}

			var accrual AccrualScheme
			if err := json.Unmarshal(data, &accrual); err != nil {
				log.Error().
					Err(err).
					Str("orderNum", order.Number).
					Msg("unmarshal accrual service response data")
				result.Error = err
				resultStream <- result
				continue
			}
			result.Order.Status = accrual.Status
			result.Order.Accrual = accrual.Amount
			resultStream <- result
		}
	}

}

func (w *AccrualWorker) getAccrualStatus(
	ctx context.Context, orderNumber string,
) ([]byte, error) {
	var res AccrualServiceResponse
	defer w.resetRetries()
	for w.hasRetries() {
		res = w.doRequest(orderNumber)

		if res.statusCode == fiber.StatusOK {
			return res.body, nil
		}

		if res.statusCode == fiber.StatusNoContent {
			return nil, errs.ErrOrdersNotPosted
		}

		w.nextRetryInterval(res.retryAfter)
		if err := w.waitRetryInterval(ctx); err != nil {
			return nil, err
		}
	}
	return nil, res.Err()
}

func (w *AccrualWorker) waitRetryInterval(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(w.currentInterval):
	}
	return nil
}

func (w *AccrualWorker) doRequest(orderNumber string) AccrualServiceResponse {
	resp := fiber.AcquireResponse()
	statusCode, body, errs := fiber.Get(w.makeRequestURL(orderNumber)).
		SetResponse(resp).
		Bytes()
	retryAfterString := string(resp.Header.Peek("Retry-After"))
	fiber.ReleaseResponse(resp)
	retryAfterInt, _ := strconv.Atoi(retryAfterString)
	retryAfter := time.Duration(retryAfterInt) * time.Second

	var err error
	if len(errs) != 0 {
		err = errs[0]
	}
	return AccrualServiceResponse{statusCode, retryAfter, body, err}
}

func (w *AccrualWorker) hasRetries() bool {
	return len(w.retryIntervals) != 0
}

func (w *AccrualWorker) nextRetryInterval(interval time.Duration) {
	if !w.hasRetries() {
		w.currentInterval = -1
		return
	}
	if interval > 0 {
		w.currentInterval = interval
		return
	}
	w.currentInterval = w.retryIntervals[0]
	w.retryIntervals = w.retryIntervals[1:]
}

func (w *AccrualWorker) resetRetries() {
	w.retryIntervals = retryIntervals[:]
}

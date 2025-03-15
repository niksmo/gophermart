package orders

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/pkg/logger"
)

var retryIntervals = [3]time.Duration{time.Second, 3 * time.Second, 5 * time.Second}

type AccrualWorkerPool struct {
	Num     int
	ChanIN  <-chan OrderScheme
	ChanOUT chan<- AccrualResult
}

func (wp AccrualWorkerPool) Run(ctx context.Context) {
	for range wp.Num {
		worker := newAccrualWorker()
		go worker.Run(ctx, wp.ChanIN, wp.ChanOUT)
	}
}

type AccrualResult struct {
	Error error
	Order OrderScheme
}

type AccrualWorker struct {
	retryIntervals []time.Duration
	makeRequestURL func(string) string
}

func newAccrualWorker() AccrualWorker {
	return AccrualWorker{
		retryIntervals: retryIntervals[:],
		makeRequestURL: config.Accrual.GetOrdersReqURL,
	}
}

func (w *AccrualWorker) Run(
	ctx context.Context,
	orderStream <-chan OrderScheme,
	resultStream chan<- AccrualResult,
) {
	for order := range orderStream {
		log := logger.Instance.With().
			Str("orderNum", order.Number).
			Caller().
			Logger()
		log.Info().Msg("receive order")
		select {
		case <-ctx.Done():
			return
		default:
			result := AccrualResult{Order: order}
			data, err := w.getAccrualStatus(ctx, order.Number)
			if err != nil {
				log.Error().Err(err).Send()
				result.Error = err
				resultStream <- result
				continue
			}

			var accrual AccrualScheme
			if err := json.Unmarshal(data, &accrual); err != nil {
				log.Error().Err(err).Send()
				result.Error = err
				resultStream <- result
				continue
			}

			log.Info().Msg("send to flush stream")
			result.Order.Status = accrual.Status
			result.Order.Accrual = accrual.Amount
			resultStream <- result
		}
	}

}

func (w *AccrualWorker) getAccrualStatus(
	ctx context.Context, orderNumber string,
) ([]byte, error) {
	var (
		interval time.Duration
		err      error
		URL      = w.makeRequestURL(orderNumber)
	)
	defer w.resetRetries()

	for len(w.retryIntervals) != 0 {
		resp := fiber.AcquireResponse()
		statusCode, body, errs := fiber.Get(URL).SetResponse(resp).Bytes()
		retryAfterString := string(resp.Header.Peek("Retry-After"))
		fiber.ReleaseResponse(resp)
		if len(errs) != 0 {
			err = errs[0]
		}

		switch {
		default:
			interval = w.getRetryInterval()
		case statusCode == fiber.StatusOK:
			return body, nil
		case statusCode == fiber.StatusNoContent:
			return nil, errors.New("No content")
		case (statusCode == fiber.StatusTooManyRequests &&
			retryAfterString != ""):
			retryAfter, err := strconv.Atoi(retryAfterString)
			if err != nil {
				interval = w.getRetryInterval()
			} else {
				interval = time.Duration(retryAfter) * time.Second
			}
		case statusCode == fiber.StatusInternalServerError:
			interval = w.getRetryInterval()
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}

	}
	return nil, err
}

func (w *AccrualWorker) getRetryInterval() time.Duration {
	interval := w.retryIntervals[0]
	w.retryIntervals = w.retryIntervals[1:]
	return interval
}

func (w *AccrualWorker) resetRetries() {
	w.retryIntervals = retryIntervals[:]
}

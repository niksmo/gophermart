package loyalty

import (
	"strconv"
	"time"

	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/di"
	"github.com/niksmo/gophermart/pkg/luhn"
)

type BalanceScheme struct {
	ID         int32     `json:"-"`
	OwnerID    int32     `json:"-"`
	Balance    float64   `json:"current"`
	Withdraw   float64   `json:"withdraw"`
	LastUpdate time.Time `json:"-"`
}

func (b *BalanceScheme) ScanRow(row di.Row) error {
	return row.Scan(
		&b.ID, &b.OwnerID, &b.Balance, &b.Withdraw, &b.LastUpdate,
	)
}

type WithdrawRequestScheme struct {
	OrderNumber string  `json:"order"`
	Amount      float64 `json:"sum"`
}

type InvalidWithdrawScheme struct {
	OrderNumber []string `json:"order,omitempty"`
	Amount      []string `json:"sum,omitempty"`
}

func (w WithdrawRequestScheme) Validate() (result InvalidWithdrawScheme, ok bool) {
	ok = true
	number, err := strconv.Atoi(w.OrderNumber)
	if err != nil || !luhn.ValidateLuhn(number) {
		result.OrderNumber = append(
			result.OrderNumber, errs.ErrOrderInvalidNum.Error(),
		)
		ok = false
	}
	if w.Amount < 0 {
		result.Amount = append(
			result.Amount, errs.ErrLoyaltyNegativeAmount.Error(),
		)
		ok = false
	}

	return
}

type WithdrawScheme struct {
	OrderNumber string    `json:"order"`
	Amount      float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

func (w *WithdrawScheme) ScanRow(row di.Row) error {
	return row.Scan(
		&w.OrderNumber, &w.Amount, &w.ProcessedAt,
	)
}

type WithdrawalsScheme []WithdrawScheme

func (wl *WithdrawalsScheme) ScanRow(row di.Row) error {
	var withdraw WithdrawScheme
	if err := withdraw.ScanRow(row); err != nil {
		return err
	}
	*wl = append(*wl, withdraw)
	return nil
}

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

func (bs *BalanceScheme) ScanRow(row di.Row) error {
	return row.Scan(
		&bs.ID, &bs.OwnerID, &bs.Balance, &bs.Withdraw, &bs.LastUpdate,
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

func (ws WithdrawRequestScheme) Validate() (result InvalidWithdrawScheme, ok bool) {
	ok = true
	number, err := strconv.Atoi(ws.OrderNumber)
	if err != nil || !luhn.ValidateLuhn(number) {
		result.OrderNumber = append(
			result.OrderNumber, errs.ErrOrderInvalidNum.Error(),
		)
		ok = false
	}
	if ws.Amount < 0 {
		result.Amount = append(
			result.Amount, errs.ErrLoyaltyNegativeAmount.Error(),
		)
		ok = false
	}

	return
}

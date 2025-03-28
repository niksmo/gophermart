package orders

import (
	"strconv"
	"time"

	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/di"
	"github.com/niksmo/gophermart/pkg/luhn"
)

type OrderNumberScheme []byte

func (orderNumber OrderNumberScheme) Validate() (string, error) {
	numberString := string(orderNumber)
	number, err := strconv.Atoi(numberString)
	if err != nil {
		return "", errs.ErrOrderInvalidFormat
	}
	if !luhn.ValidateLuhn(number) {
		return "", errs.ErrOrderInvalidNum
	}
	return numberString, nil
}

type OrderScheme struct {
	ID         int32     `json:"-"`
	OwnerID    int32     `json:"-"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float32   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (order *OrderScheme) ScanRow(row di.Row) error {
	return row.Scan(
		&order.ID,
		&order.OwnerID,
		&order.Number,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
	)
}

type OrderListScheme []OrderScheme

func (orderList *OrderListScheme) ScanRow(row di.Row) error {
	var order OrderScheme
	if err := order.ScanRow(row); err != nil {
		return err
	}
	*orderList = append(*orderList, order)
	return nil
}

type AccrualScheme struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Amount      float32 `json:"accrual,omitempty"`
}

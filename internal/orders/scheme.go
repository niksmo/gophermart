package orders

import (
	"strconv"
	"time"

	"github.com/niksmo/gophermart/internal/errs"
)

type OrderNumberScheme []byte

func (orderNumber OrderNumberScheme) Validate() (number int64, err error) {
	number, err = strconv.ParseInt(string(orderNumber), 10, 64)
	if err != nil {
		return -1, errs.ErrOrderInvalidFormat
	}
	if !isValidLuhn(number) {
		return -1, errs.ErrOrderInvalidNum
	}
	return number, nil
}

func isValidLuhn(number int64) bool {
	if number < 0 {
		return false
	}

	var sum int64
	var isEvenOrder bool
	for number > 0 {
		cur := number % 10

		if isEvenOrder {
			cur *= 2
			cur = cur/10 + cur%10
		}

		sum += cur
		number /= 10
		isEvenOrder = !isEvenOrder
	}
	return sum%10 == 0
}

type OrderScheme struct {
	ID         int32     `json:"-"`
	OwnerID    int32     `json:"-"`
	Number     int64     `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

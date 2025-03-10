package orders

import (
	"strconv"
	"time"

	"github.com/niksmo/gophermart/internal/errs"
)

type OrderNumberScheme []byte

func (orderNumber OrderNumberScheme) Validate() (string, error) {
	numberString := string(orderNumber)
	number, err := strconv.Atoi(numberString)
	if err != nil {
		return "", errs.ErrOrderInvalidFormat
	}
	if !isValidLuhn(number) {
		return "", errs.ErrOrderInvalidNum
	}
	return numberString, nil
}

func isValidLuhn(number int) bool {
	if number < 0 {
		return false
	}

	var sum int
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
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

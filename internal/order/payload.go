package order

import (
	"strconv"

	"github.com/niksmo/gophermart/internal/errs"
)

type UploadOrderReqPayload string

func (payload UploadOrderReqPayload) Validate() (number int64, err error) {
	number, err = strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		return -1, errs.ErrInvalidOrderNum
	}
	if !isValidLuhn(number) {
		return -1, errs.ErrInvalidOrderNum
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

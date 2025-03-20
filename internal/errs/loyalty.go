package errs

import "errors"

var (
	ErrLoyaltyNotEnoughPoints = errors.New("not enough bonus points")
	ErrLoyaltyNegativeAmount  = errors.New("negative amount")
	ErrLoyaltyNoWithdrawals   = errors.New("there are no points withdrawals")
)

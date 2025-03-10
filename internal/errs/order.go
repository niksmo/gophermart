package errs

import "errors"

var (
	ErrInvalidOrderNum      = errors.New("invalid order number")
	ErrOrderUploaded        = errors.New("already uploaded")
	ErrOrderUploadedByOther = errors.New("already uploaded by other user")
	ErrOrderUploadedByUser  = errors.New("you have already uploaded the order")
)

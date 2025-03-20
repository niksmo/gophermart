package errs

import "errors"

var (
	ErrOrderInvalidFormat      = errors.New("invalid pyaload data format")
	ErrOrderInvalidNum         = errors.New("invalid order number")
	ErrOrderUploaded           = errors.New("already uploaded")
	ErrOrderUploadedByOther    = errors.New("already uploaded by other user")
	ErrOrderUploadedByUser     = errors.New("already uploaded")
	ErrOrdersNoUploads         = errors.New("there are no uploaded orders")
	ErrOrdersNotPosted         = errors.New("order didn't posted to accrual system")
	ErrOrdersAccrualServiceAPI = errors.New("accrual service response with 4xx-5xx status")
)

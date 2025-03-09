package order

import "strconv"

type UploadOrderReqPayload string

func (payload UploadOrderReqPayload) Validate() (number int64, err error) {
	number, err = strconv.ParseInt(string(payload), 10, 64)
	if err != nil {
		return -1, err
	}
	return number, nil
}

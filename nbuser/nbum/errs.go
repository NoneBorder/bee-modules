package nbum

import "errors"

var (
	ErrGeneralServerError    = errors.New("service error, please try again")
	ErrUnkownError           = errors.New("got unkown error, please concat admin")
)

const (
	ErrCodeBadRequest   = 400
	ErrCodeNotLogin     = 401
	ErrCodeGeneralError = 507
	ErrCodeUnkownError  = 599
)

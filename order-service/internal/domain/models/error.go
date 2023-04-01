package domain

import "errors"

var (
	//	ErrNonExistentId       = errors.New("non-existent id")
	//	ErrIncorrectParams     = errors.New("incorrect params")
	//	ErrFailedToken         = errors.New("failed token")
	//	ErrConflict            = errors.New("conflict")
	//	ErrInternalError       = errors.New("internal server error")
	//	ErrUnauthorizedUser    = errors.New("unauthorized user")
	//	ErrContentNotFound     = errors.New("content not found")
	//	ErrOrderNotFound       = errors.New("order not found")
	ErrIdempotencyConflict = errors.New("idempotency conflict")
)

package shareaderrors

import "errors"

var (
	ErrNotFound       = errors.New("item not found")
	ErrInvalidPayload = errors.New("id and message are required")
)

package apierror

import (
	"fmt"
)

type ErrorCode string

const (
    ErrInvalidJSON ErrorCode = "invalid_json"
    ErrS3Access    ErrorCode = "s3_access_failed"
)

type APIError struct {
	Code       ErrorCode
	HTTPStatus int
	Err        error
	Internal   string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %v", e.Code, e.Err)
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, status int, err error, internal ...string) *APIError {
	apiErr := &APIError{
		Code:       code,
		HTTPStatus: status,
		Err:        err,
	}
	if len(internal) > 0 {
		apiErr.Internal = internal[0]
	}
	return apiErr
}

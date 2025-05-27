package cchttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/sts-solutions/base-code/ccvalidation"
)

// ErrorResponse represents an error object in the HTTP response
// swagger:model
type ErrorResponse struct {
	// Message is the error message
	Message string `json:"message"`
	// Code is the error code
	Code     int   `json:"code"`
	HTTPCode int   `json:"-"`
	InnerErr error `json:"-"`
}

// Error returns a string representation of the error.
func (e ErrorResponse) Error() string {
	return e.Message
}

// FrontError returns an error HTTP code (depending on error type) and response body
// - EndRequest (408)
// - RequestTimeout (408)
// - InternalServerError (500)
func FrontError(err error) (errReponse *ErrorResponse) {
	if errors.Is(err, ccvalidation.Result{}) {
		return BadRequest(err)
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return requestTimeout(err)
	}

	return InternalServerError(err)
}

// NotFound returns a NotFound (404) HTTP code and response body.
func NotFound(err error) *ErrorResponse {
	return getErrorResponse(http.StatusNotFound, err)
}

// ServiceUnavailable returns a ServiceUnavailable (503) HTTP code and response body.
func ServiceUnavailable(err error) *ErrorResponse {
	return getErrorResponse(http.StatusServiceUnavailable, err)
}

// EndRequest returns a RequestTimeout (408) HTTP code and response body
func BadRequest(err error) *ErrorResponse {
	return getErrorResponse(http.StatusBadRequest, err)
}

// RequestTimeout returns a RequestTimeout (408) HTTP code and response body
func requestTimeout(err error) *ErrorResponse {
	return getErrorResponse(http.StatusRequestTimeout, err)
}

// InternalServerError returns an InternalServerError (500) HTTP code and response body
func InternalServerError(err error) *ErrorResponse {
	return getErrorResponse(http.StatusInternalServerError, err)
}

// GetErrorResponseWithCode returns an ErrorResponse with code.
func GetErrorResponseWithCode(httpCode int, code int, err error) *ErrorResponse {
	errorReponse := ErrorResponse{
		Message:  fmt.Sprintf("%v", err),
		Code:     code,
		InnerErr: err,
		HTTPCode: httpCode,
	}

	return &errorReponse
}

func getErrorResponse(httpCode int, err error) *ErrorResponse {
	errorResponse := ErrorResponse{
		Message:  fmt.Sprintf("%v", err),
		InnerErr: err,
		HTTPCode: httpCode,
	}

	return &errorResponse
}

// GetInnerErr returns the inner error
func (e ErrorResponse) GetInnerErr() error {
	return e.InnerErr
}

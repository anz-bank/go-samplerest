package pet

import (
	"fmt"
	"net/http"
)

// Error codes
const (
	// ErrUnknown is used when an unknown error occurred
	ErrUnknown int = iota
	// ErrBadRequest is used when the incoming request is invalid
	ErrBadRequest
	// ErrInvalidID is used when an invalid ID is encountered
	ErrInvalidID
	// ErrIDAlreadyExists is used when attempting to erroneously overwrite an existing entry
	ErrIDAlreadyExists
	// ErrIDNotFound is used when attempting to read a non-existing entry
	ErrIDNotFound
)

// APIError defines an error that separates internal and external error messages
type APIError struct {
	Message string
	Code    int
	Cause   error
}

func (e *APIError) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%v : %v", e.Message, e.Cause)
}

// APIErrorf creates a new ApiError with formatting
func APIErrorf(code int, format string, args ...interface{}) *APIError {
	return APIErrorEf(code, nil, format, args...)
}

// APIErrorEf creates a new APIError with causing error and formatting
func APIErrorEf(code int, cause error, format string, args ...interface{}) *APIError {
	return &APIError{
		Message: fmt.Sprintf(format, args...),
		Code:    code,
		Cause:   cause,
	}
}

var errStatusMap = map[int]int{
	ErrBadRequest:      http.StatusBadRequest,
	ErrUnknown:         http.StatusInternalServerError,
	ErrIDAlreadyExists: http.StatusInternalServerError,
	ErrIDNotFound:      http.StatusNotFound,
}

func renderHTTPErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	if apiError, ok := err.(*APIError); ok {
		statusCode, ok := errStatusMap[apiError.Code]
		if ok {
			http.Error(w, apiError.Message, statusCode)
			return
		}
		http.Error(w, apiError.Message, http.StatusInternalServerError)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return
}

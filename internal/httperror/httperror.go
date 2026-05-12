// https://medium.com/@alexisbouchez/centralize-http-error-handling-in-go-ba78446f4d9d

package utility

import (
	"errors"
	"net/http"
)

type HTTPError struct {
	error
	Code int
}

func New(code int, message string) *HTTPError {
	return &HTTPError{
		error: errors.New(message),
		Code:  code,
	}
}
func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, message)
}

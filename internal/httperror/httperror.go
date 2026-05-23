// https://medium.com/@alexisbouchez/centralize-http-error-handling-in-go-ba78446f4d9d

package httperror

import (
	"errors"
	"net/http"
)

type HTTPError struct {
	error
	Code int
	Err  error
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

func New(code int, message string, err error) *HTTPError {
	return &HTTPError{
		error: errors.New(message),
		Code:  code,
		Err:   err,
	}
}

func NotFound(message string) *HTTPError {
	return New(http.StatusNotFound, message, nil)
}

func BadRequest(message string) *HTTPError {
	return New(http.StatusBadRequest, message, nil)
}

func InternalServerError(err error) *HTTPError {
	return New(http.StatusInternalServerError, "Something catastrophic has happened", err)
}

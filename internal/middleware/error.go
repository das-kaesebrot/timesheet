package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/das-kaesebrot/timesheet/internal/httperror"
	"github.com/das-kaesebrot/timesheet/internal/template"
)

func ErrorHandler(renderer *template.Renderer) func(func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(next func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			err := next(w, r)
			if err != nil {
				var httpErr *httperror.HTTPError
				var code int
				var message string

				if errors.As(err, &httpErr) {
					log.Printf("HTTP %d %s (cause: %v)", httpErr.Code, httpErr.Error(), httpErr.Unwrap())
					w.WriteHeader(httpErr.Code)
					code = httpErr.Code
					message = httpErr.Error()
				} else {
					log.Printf("unhandled error: %v", err)
					code = http.StatusInternalServerError
					message = "Internal Server Error"
				}

				w.WriteHeader(code)
				renderer.Render(w, "error", map[string]interface{}{
					"Code":    code,
					"Message": message,
				})
			}
		}
	}
}

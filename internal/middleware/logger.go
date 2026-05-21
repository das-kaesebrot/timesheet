package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() { log.Println(r.URL.Path, time.Now()) }()
		next(w, r)
	}
}

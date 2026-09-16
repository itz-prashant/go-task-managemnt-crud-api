package middleware

import (
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Panicf("[HTTP %s %s from %s took %w", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	})
}

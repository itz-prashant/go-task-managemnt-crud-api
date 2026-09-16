package middleware

import (
	"log"
	"net/http"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer func(){
			if err := recover(); err != nil {
				log.Printf("[PANIC RECOVERED] %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error}:"internal server eror"`))
			}
		}()
		next.ServeHTTP(w,r)
	})
}
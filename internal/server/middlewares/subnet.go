package middlewares

import (
	"log"
	"net/http"
)

func SubnetValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checiks subnet")
		next.ServeHTTP(w, r)
	})
}

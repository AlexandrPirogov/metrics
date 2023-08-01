package middlewares

import (
	"log"
	"net/http"
)

func SubnetValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request came with IP-FROM %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

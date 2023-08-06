package middlewares

import (
	"memtracker/internal/server/verifier"
	"net/http"
)

func SubnetValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok := verifier.IsContainerSubnetHTTP(r); !ok {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

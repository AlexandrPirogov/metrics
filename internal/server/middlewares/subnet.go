package middlewares

import (
	"log"
	"memtracker/internal/config/server"
	"net/http"
	"net/netip"
)

func SubnetValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cameIP := r.Header.Get("X-Real-IP")
		ip, err := netip.ParseAddr(cameIP)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		subnet, _ := netip.ParsePrefix(server.ServerCfg.Subnet)
		if subnet.Contains(ip) {
			next.ServeHTTP(w, r)
			return
		}

		w.WriteHeader(http.StatusForbidden)
	})
}

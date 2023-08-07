package verifier

import (
	"log"
	"memtracker/internal/config/server"
	"net/http"
	"net/netip"

	"google.golang.org/grpc/metadata"
)

const xRealIP = "X-Real-IP"

func IsContainerSubnetHTTP(r *http.Request) bool {
	cameIP := r.Header.Get(xRealIP)
	ip, parseErr := netip.ParseAddr(cameIP)
	subnet, parsePrefixErr := netip.ParsePrefix(server.ServerCfg.Subnet)
	return parseErr == nil && parsePrefixErr == nil && subnet.Contains(ip)
}

func IsContainerSubnetGRPC(md metadata.MD) bool {
	if server.ServerCfg.Subnet == "" {
		return true
	}

	cameIP := md.Get(xRealIP)
	if len(cameIP) == 0 {
		return false
	}

	ip, parseErr := netip.ParseAddr(cameIP[0])

	subnet, parsePrefixErr := netip.ParsePrefix(server.ServerCfg.Subnet)
	log.Println(ip)
	return parseErr == nil && parsePrefixErr == nil && subnet.Contains(ip)
}

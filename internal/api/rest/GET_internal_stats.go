package restapi

import (
	"net"
	"net/http"

	"github.com/aube/url-shortener/internal/app/config"
	"github.com/aube/url-shortener/internal/app/usecases"
	"github.com/aube/url-shortener/internal/logger"
)

var usecasesGetStats = usecases.GetStats

// HandlerInternalStats creates an HTTP handler for retrieving internal statistics.
// The handler checks the requesting IP against a trusted subnet before providing access.
// Only requests from the trusted subnet (configured in config.TrustedSubnet) are allowed.
// The client IP is expected in the X-Real-IP header.
// Returns:
//
//	http.HandlerFunc - the HTTP handler function
func HandlerInternalStats() http.HandlerFunc {
	config := config.NewConfig()

	return func(w http.ResponseWriter, r *http.Request) {

		if config.TrustedSubnet == "" {
			http.Error(w, "wrong IP", http.StatusForbidden)
			return
		}

		ip := r.Header.Get("X-Real-IP")
		isTrustedAddress := isIPInCIDR(ip, config.TrustedSubnet)

		if !isTrustedAddress {
			http.Error(w, "wrong IP", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		log := logger.WithContext(ctx)

		json, err := usecasesGetStats(ctx)

		if err != nil {
			http.Error(w, "Error on stats reuest", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)

		n, err := w.Write(json)
		if err != nil {
			// Handle error (connection may have been closed)
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}

		log.Info("HandlerInternalStats", "Wrote bytes", n, "json", json)
	}
}

// isIPInCIDR checks if an IP address is within a specified CIDR range.
// Parameters:
//
//	ipStr - the IP address to check (string format)
//	cidrStr - the CIDR range to check against (string format)
//
// Returns:
//
//	bool - true if the IP is within the CIDR range, false otherwise
func isIPInCIDR(ipStr, cidrStr string) bool {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	return ipnet.Contains(ip)
}

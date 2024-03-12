package middleware

import (
	"net/http"
	"net/netip"

	"github.com/kuzhukin/metrics-collector/internal/zlog"
)

type IPChecker struct {
	trustedSubnet netip.Prefix
}

func NewIPChecker(subnet netip.Prefix) *IPChecker {
	return &IPChecker{trustedSubnet: subnet}
}

func (c *IPChecker) CheckIPHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIP := r.Header.Get("X-Real-IP")

		ip, err := netip.ParseAddr(realIP)
		if err != nil {
			zlog.Logger.Errorf("parse ip addr err=%s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !c.trustedSubnet.Contains(ip) {
			zlog.Logger.Errorf("no such IP=%s in trusted subnet=%s", ip, c.trustedSubnet)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

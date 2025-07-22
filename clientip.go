package http

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// ClientIP returns the client's IP address from HTTP headers
func ClientIP(req *http.Request) (string, error) {
	var (
		err  error
		ip   net.IP
		host string
	)

	if req == nil {
		return "", errors.New("*http.Request must not be nil")
	}

	// Headers to check in priority order
	headers := []string{
		"True-Client-IP",   // Cloudflare Enterprise
		"CF-Connecting-IP", // Cloudflare
		"Fastly-Client-IP", // Fastly CDN
		"Fly-Client-IP",    // Fly.io
		"X-Real-IP",        // Common reverse proxies
		"X-Client-IP",      // Alternative common header
		"X-Forwarded-For",  // Standard proxy header
	}

	for _, header := range headers {
		if values := req.Header.Values(header); len(values) > 0 {
			for _, value := range values {
				ips := strings.Split(value, ",")
				for _, ipStr := range ips {
					ipStr = strings.TrimSpace(ipStr)
					if ip = net.ParseIP(ipStr); ip != nil {
						return ipStr, nil
					}
				}
			}
		}
	}

	// Fallback to direct connection IP
	if host, _, err = net.SplitHostPort(req.RemoteAddr); err != nil {
		return "", err
	}
	if net.ParseIP(host) == nil {
		return "", errors.New(fmt.Sprintf("invalid remote address format %s", ip))
	}

	return host, nil
}

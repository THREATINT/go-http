package http

import (
	"errors"
	"net"
	"net/http"
	"strings"
)

// GetClientIP (request)
// Returns the client ip address from HTTP headers
func GetClientIP(r *http.Request) (string, error) {
	var (
		err           error
		xForwardedFor string
		ip            string
		clientIP      net.IP
	)

	if r == nil {
		return "", errors.New("*http.Request must not be nil")
	}

	// True-Client-IP, used by CloudFlare (attn: Enterprise plan only!)
	if cfTrueClientIP := r.Header.Get("True-Client-IP"); cfTrueClientIP != "" {
		return cfTrueClientIP, nil
	}
	// CF-Connecting-IP, used by CloudFlare
	if cfConnectingIP := r.Header.Get("CF-Connecting-IP"); cfConnectingIP != "" {
		return cfConnectingIP, nil
	}

	// used by most reverse proxies
	if xForwardedFor = r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
		// we are only interested in the client-ip
		// (see https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Forwarded-For)
		c := strings.SplitN(xForwardedFor, ",", 2)
		if len(c) == 2 {
			return c[0], nil
		}

		return xForwardedFor, nil
	}

	if ip, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
		return "", errors.New("cannot parse remote address ( ip + port)")
	}

	if clientIP = net.ParseIP(ip); clientIP == nil {
		return "", errors.New("cannot parse remote ip")
	}

	return ip, nil
}

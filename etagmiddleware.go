package http

import (
	"bytes"
	"crypto/sha3"
	"fmt"
	"hash"
	"io"
	"net/http"
)

func EtagMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Only apply to HTTP GET and HEAD, see RFC #9110
		if req.Method != http.MethodGet && req.Method != http.MethodHead {
			next.ServeHTTP(rw, req)
			return
		}

		etagWriter := &EtagResponseWriter{
			ResponseWriter: rw,
			hash:           sha3.New224(),
			buf:            bytes.Buffer{},
		}
		etagWriter.w = io.MultiWriter(&etagWriter.buf, etagWriter.hash)

		next.ServeHTTP(etagWriter, req)

		finalStatus := etagWriter.status
		if finalStatus == 0 {
			finalStatus = http.StatusOK
		}

		// Use ETag only for successful responses, see RFC #9110
		if finalStatus == http.StatusOK {
			computedETag := fmt.Sprintf("\"%x\"", etagWriter.hash.Sum(nil))

			if clientETag := req.Header.Get("If-None-Match"); clientETag == computedETag {
				rw.Header().Set("ETag", computedETag)
				rw.WriteHeader(http.StatusNotModified)
				return
			}

			rw.Header().Set("ETag", computedETag)
		}

		rw.WriteHeader(finalStatus)

		if _, err := etagWriter.buf.WriteTo(rw); err != nil {
			panic(err)
		}
	})
}

type EtagResponseWriter struct {
	http.ResponseWriter
	buf    bytes.Buffer
	w      io.Writer
	hash   hash.Hash
	status int
}

func (e *EtagResponseWriter) Header() http.Header {
	return e.ResponseWriter.Header()
}

func (e *EtagResponseWriter) Write(p []byte) (int, error) {
	if e.status == 0 {
		e.status = http.StatusOK
	}
	return e.w.Write(p)
}

func (e *EtagResponseWriter) WriteHeader(statusCode int) {
	if e.status == 0 {
		e.status = statusCode
	}
}

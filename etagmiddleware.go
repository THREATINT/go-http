package http

import (
	"bytes"
	"crypto/sha3"
	"fmt"
	"hash"
	"io"
	"net/http"
)

// ETagMiddleware to create an ETag according to https://www.rfc-editor.org/rfc/rfc9110#name-etag
func ETagMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(respW http.ResponseWriter, req *http.Request) {
		etagWriter := &EtagResponseWriter{
			ResponseWriter: respW,
			hash:           sha3.New224(),
			buf:            bytes.Buffer{},
		}
		etagWriter.w = io.MultiWriter(&etagWriter.buf, etagWriter.hash)

		next.ServeHTTP(etagWriter, req)

		computedETag := fmt.Sprintf("\"%x\"", etagWriter.hash.Sum(nil))

		for k, v := range etagWriter.ResponseWriter.Header() {
			respW.Header()[k] = v
		}
		respW.Header().Set("ETag", computedETag)

		if clientETag := req.Header.Get("If-None-Match"); clientETag == computedETag {
			respW.WriteHeader(http.StatusNotModified)
			return
		}

		if etagWriter.status == 0 {
			respW.WriteHeader(http.StatusOK)
		} else {
			respW.WriteHeader(etagWriter.status)
		}

		if _, err := etagWriter.buf.WriteTo(respW); err != nil {
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

package http

import (
	"bytes"
	"crypto/sha3"
	"fmt"
	"hash"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

// ETagMiddleware middleware to add strong ETags  ,
// see https://www.rfc-editor.org/rfc/rfc9110#name-etag for details
func ETagMiddleware(config *ETagMiddlewareConfig) func(http.Handler) http.Handler {
	if config != nil {
		for _, u := range config.IgnoreURLs {
			config.ignoreURLsRegEx = append(config.ignoreURLsRegEx, regexp.MustCompile(u))
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(respW http.ResponseWriter, req *http.Request) {
			if config.ignoreURLsRegEx != nil {
				url := req.URL.String()
				for _, rx := range config.ignoreURLsRegEx {
					if rx.Match([]byte(url)) {
						next.ServeHTTP(respW, req)
						return
					}
				}
			}

			etagWriter := &etagResponseWriter{
				ResponseWriter: respW,
				hash:           sha3.New256(),
				buf:            bytes.Buffer{},
			}
			etagWriter.w = io.MultiWriter(&etagWriter.buf, etagWriter.hash)

			next.ServeHTTP(etagWriter, req)

			// handle edge case: application logic didn't write a status code
			if etagWriter.status == 0 {
				etagWriter.status = http.StatusOK
			}

			computedETag := fmt.Sprintf("\"%x\"", etagWriter.hash.Sum(nil))

			if etagWriter.status == http.StatusOK || etagWriter.status == http.StatusCreated {
				respW.Header().Set("ETag", computedETag)
				if containsETag(req.Header.Get("If-None-Match"), computedETag) {
					if req.Method == "GET" || req.Method == "HEAD" {
						respW.WriteHeader(http.StatusNotModified)
					} else {
						respW.WriteHeader(http.StatusPreconditionFailed)
					}

					return
				}
			}

			respW.WriteHeader(etagWriter.status)

			if _, err := etagWriter.buf.WriteTo(respW); err != nil {
				if config == nil {
					panic(err)
				}
				if config.Panic {
					panic(err)
				} else {
					// respW.WriteHeader(http.StatusInternalServerError)
					// MUST NOT be used, because we have already started to write the body.
					// Writing another header would now result in a panic.

					// If the config has a Logger, we log the error message
					if config.Logger != nil {
						config.Logger.Println(err)
					}
				}
			}
		})
	}
}

// containsETag checks if the given ETag is present in the If-None-Match header value.
// If-None-Match may contain:
// - A single ETag: "etag-value"
// - Multiple ETags: "etag1", "etag2", "etag3"
// - Wildcard: *
//
// see https://www.rfc-editor.org/rfc/rfc9110.html#name-if-none-match for details
func containsETag(ifNoneMatch, computedETag string) bool {
	ifNoneMatch = strings.TrimSpace(ifNoneMatch)

	if ifNoneMatch == "" {
		return false
	}

	if ifNoneMatch == "*" {
		return true
	}

	for _, etag := range strings.Split(ifNoneMatch, ",") {
		if strings.TrimPrefix(strings.TrimSpace(etag), "W/") == computedETag {
			return true
		}
	}

	return false
}

type ETagMiddlewareConfig struct {
	Panic           bool
	Logger          *log.Logger
	IgnoreURLs      []string
	ignoreURLsRegEx []*regexp.Regexp
}

type etagResponseWriter struct {
	http.ResponseWriter
	buf    bytes.Buffer
	w      io.Writer
	hash   hash.Hash
	status int
}

func (e *etagResponseWriter) Header() http.Header {
	return e.ResponseWriter.Header()
}

func (e *etagResponseWriter) Write(p []byte) (int, error) {
	if e.status == 0 {
		e.status = http.StatusOK
	}
	return e.w.Write(p)
}

func (e *etagResponseWriter) WriteHeader(statusCode int) {
	if e.status == 0 {
		e.status = statusCode
	}
}

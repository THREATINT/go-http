package http

import (
	"mime"
	"strings"
)

// MimeTypeByExtension (extension)
//
func MimeTypeByExtension(extension string) string {
	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	m := mime.TypeByExtension(extension)

	if m != "" {
		return m
	}

	return mimeType[strings.TrimLeft(extension, ".")]
}

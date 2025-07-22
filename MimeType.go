package http

import (
	"mime"
	"strings"
)

// MimeTypeByExtension (extension)
// Return matching MimeType for file extension
func MimeTypeByExtension(extension string) string {
	var m string

	if !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}
	if m = mime.TypeByExtension(extension); m != "" {
		return m
	}

	return mimeType[strings.TrimLeft(extension, ".")]
}

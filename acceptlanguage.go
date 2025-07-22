package http

import (
	"sort"
	"strconv"
	"strings"
)

// AcceptedLanguage represents a parsed language with quality factor, see RFC #9110
type AcceptedLanguage struct {
	Lang string
	Q    float64
}

// ParseAcceptLanguage parses the Accept-Language header
func ParseAcceptLanguage(header string) []AcceptedLanguage {
	header = strings.TrimSpace(header)
	if header == "" {
		return nil
	}

	langs := make([]AcceptedLanguage, 0, 5)

	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split language tag and parameters
		segments := strings.Split(part, ";")
		langTag := strings.TrimSpace(segments[0])

		// Default quality factor: "If no "q" parameter is present, the default weight is 1.",
		// see , see https://www.rfc-editor.org/rfc/rfc9110.html#name-quality-values
		qval := 1.0

		// Process quality factor parameters
		for _, param := range segments[1:] {
			param = strings.TrimSpace(param)
			if !strings.HasPrefix(param, "q=") {
				continue // Skip non-q parameters
			}

			// Parse quality factor
			qstr := strings.TrimPrefix(param, "q=")
			qval, _ = strconv.ParseFloat(qstr, 64)

			break // Only first q= parameter matters, see RFC #9110 Section 5.3
		}

		// Only include languages with q in [ 0.001, 1.000 ], see https://www.rfc-editor.org/rfc/rfc9110.html#name-quality-values
		if qval < 0.001 || qval > 1.0 {
			continue
		}

		langs = append(langs, AcceptedLanguage{
			Lang: langTag, // Preserve case for normalization
			Q:    qval,
		})
	}

	// Order by descending quality factor, see RFC #9110 Section 12.4.4
	sort.SliceStable(langs, func(i, j int) bool {
		if langs[i].Q == langs[j].Q {
			// Maintain original header order for same quality, see RFC #9110 Section 12.4.4
			return i < j
		}
		return langs[i].Q > langs[j].Q
	})

	return langs
}

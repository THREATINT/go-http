package http

import (
	"sort"
	"strconv"
	"strings"
)

// AcceptedLanguage represents a parsed language with quality factor,
// see https://www.rfc-editor.org/rfc/rfc9110.html#name-quality-values
type AcceptedLanguage struct {
	Lang string
	Q    float64
}

// ParseAcceptLanguage parses the Accept-Language header,
// see https://www.rfc-editor.org/rfc/rfc9110.html#name-quality-values
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
		qval := 1.0

		for _, param := range segments[1:] {
			param = strings.TrimSpace(param)
			if !strings.HasPrefix(param, "q=") {
				continue
			}

			qstr := strings.TrimPrefix(param, "q=")
			qval, _ = strconv.ParseFloat(qstr, 64)

			// Only first q= parameter matters,
			// see https://www.rfc-editor.org/rfc/rfc9110.html#name-field-order
			break
		}

		// Only include languages with q in [ 0.001, 1.000 ]
		if qval < 0.001 || qval > 1.0 {
			continue
		}

		langs = append(langs, AcceptedLanguage{
			Lang: langTag,
			Q:    qval,
		})
	}

	sort.SliceStable(langs, func(i, j int) bool {
		if langs[i].Q == langs[j].Q {
			return i < j
		}
		return langs[i].Q > langs[j].Q
	})

	return langs
}

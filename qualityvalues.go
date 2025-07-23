package http

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// QualityValue represents a values with quality factor,
// see https://www.rfc-editor.org/rfc/rfc9110.html#name-quality-values
type QualityValue struct {
	Q     float64
	Value string
}

// ParseQualityValues parses a string containing values and quality factor,
// see https://www.rfc-editor.org/rfc/rfc9110.html#name-quality-values
func ParseQualityValues(header string) ([]QualityValue, error) {
	header = strings.TrimSpace(header)
	if header == "" {
		return nil, errors.New("header string must not be empty")
	}

	var (
		qualityvalues = make([]QualityValue, 0, 5)
		err           error
		q             float64 = 1.0 // default quality factor: "If no "q" parameter is present, the default weight is 1."
	)

	for _, part := range strings.Split(header, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Split language tag and parameters
		segments := strings.Split(part, ";")
		langTag := strings.TrimSpace(segments[0])

		for _, param := range segments[1:] {
			param = strings.TrimSpace(param)

			if !strings.HasPrefix(param, "q=") {
				continue
			}

			qstr := strings.TrimPrefix(param, "q=")
			if q, err = strconv.ParseFloat(qstr, 64); err != nil {
				return nil, errors.New(fmt.Sprintf("invalid value 'q=%s'", qstr))
			}

			// Only first q= parameter matters,
			// see https://www.rfc-editor.org/rfc/rfc9110.html#name-field-order
			break
		}

		// Only include languages with q in [ 0.001, 1.000 ]
		if q < 0.001 || q > 1.0 {
			continue
		}

		qualityvalues = append(qualityvalues, QualityValue{
			Value: langTag,
			Q:     q,
		})
	}

	sort.SliceStable(qualityvalues, func(i, j int) bool {
		if qualityvalues[i].Q == qualityvalues[j].Q {
			return i < j
		}
		return qualityvalues[i].Q > qualityvalues[j].Q
	})

	return qualityvalues, nil
}

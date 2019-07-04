package http

import (
	"strconv"
	"strings"
)

type AcceptedLanguage struct {
	Lang string
	Q    float64
}

func ParseAcceptLanguage(acceptLanguage string) []AcceptedLanguage {
	var acceptedLanguages []AcceptedLanguage

	for _, langs := range strings.Split(acceptLanguage, ",") {
		langQ := strings.Split(strings.Trim(langs, " "), ";")
		if len(langQ) == 1 {
			acceptedLanguages = append(acceptedLanguages, AcceptedLanguage{Lang: strings.ToLower(langQ[0]), Q: 1.0})
		} else {
			q, err := strconv.ParseFloat(strings.Split(langQ[1], "=")[1], 64)
			if err != nil {
				q = 0.0
			}
			acceptedLanguages = append(acceptedLanguages, AcceptedLanguage{Lang: strings.ToLower(langQ[0]), Q: q})
		}
	}

	return acceptedLanguages
}

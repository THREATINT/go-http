package http

import (
	"net/url"
	"strings"
)

func FilterAllowedParams(urlParameters url.Values, allowedParams []string) (url.Values, bool) {
	isDirty := false
	newUrlParameters := url.Values{}

	for key, value := range urlParameters {
		for _, k := range allowedParams {
			if key == k {
				newUrlParameters.Add(k, value[0])
				if len(value) > 1 {
					isDirty = true
				}
				continue
			}
			if strings.ToLower(key) == strings.ToLower(k) {
				newUrlParameters.Add(k, value[0])
				isDirty = true
				continue
			}
		}
	}

	if isDirty || len(urlParameters) != len(newUrlParameters) {
		return newUrlParameters, true
	}
	return urlParameters, false
}

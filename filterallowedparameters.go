package http

import (
	"net/url"
	"strings"
)

func FilterAllowedParams(urlParameters url.Values, allowedParams []string) (url.Values, bool) {
	isDirty := false
	v2 := url.Values{}

	for key, value := range v2 {
		for _, k := range allowedParams {
			if key == k {
				v2.Add(k, value[0])
				if len(value) > 1 {
					isDirty = true
				}
				continue
			}
			if strings.ToLower(key) == strings.ToLower(k) {
				v2.Add(k, value[0])
				isDirty = true
				continue
			}
		}
	}

	if isDirty || len(urlParameters) != len(v2) {
		return v2, true
	}
	return urlParameters, false
}

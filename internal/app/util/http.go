package util

import (
	"net/http"
	"strings"
)

func IsFilteredWithHttp(r *http.Request) bool {
	var path = r.RequestURI

	for _, word := range strings.Split(path, "/") {
		if word == "health" || word == "debug" || word == "metric" {
			return true
		}
	}

	return false
}

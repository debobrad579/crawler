package main

import (
	"net/url"
	"strings"
)

func normalizeURL(urlString string) (string, error) {
	urlStruct, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	return urlStruct.Host + strings.TrimSuffix(urlStruct.Path, "/"), nil
}

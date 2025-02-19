package main

import (
	"net/url"
)

func normalizeURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	nURL := parsed.Host + parsed.Path
	return nURL, nil
}

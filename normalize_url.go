package main

import (
	"errors"
	"net/url"
	"path"
)

func normalizeURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsed.Host == "" {
		return "", errors.New("invalid host")
	}
	nURL := parsed.Host
	if parsed.Path != "" {
		nURL += path.Clean(parsed.Path)
	}
	return nURL, nil
}

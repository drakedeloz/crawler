package main

import (
	"errors"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/html"
)

func normalizeURL(rawURL string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	if parsed.Host == "" {
		return "", errors.New("invalid host")
	}
	nURL := parsed.Scheme + "://" + parsed.Host
	if parsed.Path != "" {
		nURL += path.Clean(parsed.Path)
	}
	return nURL, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return []string{}, err
	}
	htmlReader := strings.NewReader(htmlBody)
	nodeTree, err := html.Parse(htmlReader)
	if err != nil {
		return []string{}, err
	}
	anchorNodes := searchNodes(nodeTree)
	var urls []string
	for _, node := range anchorNodes {
		href := getHref(&node)
		href = validateURL(baseURL, href)
		urls = append(urls, href)
	}
	return urls, nil
}

func searchNodes(node *html.Node) []html.Node {
	anchorNodes := []html.Node{}
	if node.Type == html.ElementNode && node.Data == "a" {
		anchorNodes = append(anchorNodes, *node)
	}
	if node.FirstChild != nil {
		anchorNodes = append(anchorNodes, searchNodes(node.FirstChild)...)
	}
	if node.NextSibling != nil {
		anchorNodes = append(anchorNodes, searchNodes(node.NextSibling)...)
	}
	return anchorNodes
}

func getHref(anchorNode *html.Node) string {
	for _, attr := range anchorNode.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func validateURL(baseURL *url.URL, href string) string {
	parsedHref, err := url.Parse(href)
	if err != nil {
		return ""
	}
	if !parsedHref.IsAbs() {
		parsedHref = baseURL.ResolveReference(parsedHref)
	}
	return parsedHref.String()
}

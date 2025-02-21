package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}
	pages := make(map[string]int)
	baseURL := args[0]
	fmt.Printf("starting crawl of: %s\n", baseURL)
	crawlPage(baseURL, baseURL, pages)
}

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	if resp.StatusCode > 399 {
		return "", errors.New("failed to get html")
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", errors.New("invalid content type")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	if !sameDomain(rawBaseURL, rawCurrentURL) {
		return
	}
	nCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, found := pages[nCurrent]; !found {
		pages[nCurrent] = 1
	} else {
		pages[nCurrent]++
		return
	}
	fmt.Printf("Getting HTML for %s\n", nCurrent)
	rawHTML, err := getHTML(nCurrent)
	if err != nil {
		fmt.Println(err)
		return
	}
	foundURLs, err := getURLsFromHTML(rawHTML, rawBaseURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, url := range foundURLs {
		crawlPage(rawBaseURL, url, pages)
	}
}

func sameDomain(rawBaseURL, rawCurrentURL string) bool {
	nBase, err := normalizeURL(rawBaseURL)
	if err != nil {
		fmt.Println("failed to normalize rawbase")
		return false
	}
	nCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("failed to normalize rawcurrent")
		return false
	}
	pBase, err := url.Parse(nBase)
	if err != nil {
		fmt.Println(err)
		return false
	}
	pCurrent, err := url.Parse(nCurrent)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return pBase.Host == pCurrent.Host
}

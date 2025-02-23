package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	maxPages           int
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("Usage: crawler <website> <max_concurrency> <max_pages>")
		os.Exit(1)
	}

	maxConcurrency, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Invalid max concurrency value: %v\n", args[1])
		os.Exit(1)
	}

	maxPage, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("Invalid max pages value: %v\n", args[2])
		os.Exit(1)
	}

	nBase, err := url.Parse(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	cfg := &config{
		pages:              make(map[string]int),
		baseURL:            nBase,
		maxPages:           maxPage,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
	}
	cfg.wg.Add(1)
	cfg.crawlPage(cfg.baseURL.String())
	cfg.wg.Wait()
	printReport(cfg.pages, cfg.baseURL.String())
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

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer cfg.wg.Done()
	defer func() { <-cfg.concurrencyControl }()
	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()
	if !sameDomain(cfg.baseURL.String(), rawCurrentURL) {
		return
	}
	nCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	if isFirst := cfg.addPageVisit(nCurrent); isFirst {
		rawHTML, err := getHTML(nCurrent)
		if err != nil {
			fmt.Println(err)
			return
		}
		foundURLs, err := getURLsFromHTML(rawHTML, cfg.baseURL.String())
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, url := range foundURLs {
			cfg.wg.Add(1)
			go cfg.crawlPage(url)
		}
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

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, found := cfg.pages[normalizedURL]; !found {
		cfg.pages[normalizedURL] = 1
		return true
	} else {
		cfg.pages[normalizedURL]++
		return false
	}
}

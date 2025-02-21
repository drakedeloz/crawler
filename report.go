package main

import (
	"fmt"
	"sort"
)

type Page struct {
	url   string
	count int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("=============================")
	sortedPages := sortMap(pages)
	for _, page := range sortedPages {
		fmt.Printf("Found %d internal links to %s\n", page.count, page.url)
	}
}

func sortMap(pages map[string]int) []Page {
	var sortedPages []Page
	for key, val := range pages {
		sortedPages = append(sortedPages, Page{key, val})
	}
	sort.Slice(sortedPages, func(i, j int) bool {
		if sortedPages[i].count == sortedPages[j].count {
			return sortedPages[i].url < sortedPages[j].url
		}
		return sortedPages[i].count > sortedPages[j].count
	})
	return sortedPages
}

# Crawler

Crawler is a simple CLI web crawler that generates a report for internal links found on a given domain.

You can limit the number of pages to crawl and control the maximum concurrency of the crawler process.

### Usage

Arguments:

- ```<maxConcurrency>```: The maximum number of goroutines (threads) the crawler should use to fetch pages. (Type: int)
- ```<maxPages>```: The maximum number of pages to crawl. (Type: int)

Example Command:

```./crawler 10 50```

This command sets the maximum concurrency to 10 and limits the crawler to 50 pages.

package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

func getHTML(rawURL string) (string, error) {
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "BootCrawler/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("http error: %v", resp.Status)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "text/html") {
		return "", fmt.Errorf("invalid content type: %s", contentType)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

type config struct {
	maxPages           int
	pages              map[string]PageData
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	defer cfg.wg.Done()

	cfg.concurrencyControl <- struct{}{}
	defer func() { <-cfg.concurrencyControl }()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil || currentURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}
	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if !cfg.addPageVisit(normalizedCurrentURL) {
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("crawled url: %s\n", rawCurrentURL)

	data, err := extractPageData(html, rawCurrentURL)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, outgoingLink := range data.OutgoingLinks {
		cfg.wg.Add(1)
		go cfg.crawlPage(outgoingLink)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if len(cfg.pages) >= cfg.maxPages {
		return false
	}

	if _, ok := cfg.pages[normalizedURL]; !ok {
		cfg.pages[normalizedURL] = PageData{}
		return true
	}

	return false
}

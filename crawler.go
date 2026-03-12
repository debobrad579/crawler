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

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil || currentURL.Hostname() != baseURL.Hostname() {
		return
	}
	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	if _, ok := pages[normalizedCurrentURL]; ok {
		pages[normalizedCurrentURL]++
		return
	} else {
		pages[normalizedCurrentURL] = 1
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("an error occured fetching html: %v\n", err)
		return
	}
	fmt.Printf("fetched url: %s\n", rawCurrentURL)

	data, err := extractPageData(html, rawCurrentURL)
	if err != nil {
		fmt.Printf("an error occured extracting page data: %v\n", err)
		return
	}

	for _, outgoingLink := range data.OutgoingLinks {
		crawlPage(rawBaseURL, outgoingLink, pages)
	}
}

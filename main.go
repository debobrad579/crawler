package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	maxConcurrency, maxPages := 5, 25
	var err error

	if len(args) > 1 {
		maxConcurrency, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Println("invalid max concurrency")
			os.Exit(1)
		}
	}

	if len(args) > 2 {
		maxPages, err = strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("invalid max size")
			os.Exit(1)
		}
	}

	rawBaseURL := args[0]
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := &config{
		maxPages:           maxPages,
		baseURL:            baseURL,
		pages:              make(map[string]PageData),
		mu:                 &sync.Mutex{},
		wg:                 &sync.WaitGroup{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	writeJSONReport(cfg.pages, safeFilenameFromURL(baseURL))
}

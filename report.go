package main

import (
	"encoding/json"
	"net/url"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type PageReport struct {
	URL            string   `json:"url"`
	Heading        string   `json:"heading"`
	FirstParagraph string   `json:"first_paragraph"`
	OutgoingLinks  []string `json:"outgoing_links"`
	ImageURLs      []string `json:"image_urls"`
}

func writeJSONReport(pages map[string]PageData, filename string) error {
	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	reports := make([]PageReport, len(pages))
	for i, k := range keys {
		pageData := pages[k]
		if pageData.URL == nil {
			continue
		}

		reports[i] = PageReport{
			URL:            pageData.URL.String(),
			Heading:        pageData.Heading,
			FirstParagraph: pageData.FirstParagraph,
			OutgoingLinks:  pageData.OutgoingLinks,
			ImageURLs:      pageData.OutgoingLinks,
		}
	}

	data, err := json.MarshalIndent(reports, "", " ")
	if err != nil {
		return err
	}

	err = os.MkdirAll("reports", 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join("reports", filename), data, 0666)
	return err
}

func safeFilenameFromURL(urlStruct *url.URL) string {
	filename := urlStruct.Host + urlStruct.Path

	re := regexp.MustCompile(`[<>:"/\\|?*]`)
	filename = re.ReplaceAllString(filename, "_")

	filename = strings.Trim(filename, "_")

	return filename
}

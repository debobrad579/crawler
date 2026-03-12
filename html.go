package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	h1 := doc.Find("h1")
	if h1.Text() != "" {
		return h1.Text()
	}

	return doc.Find("h2").Text()
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	main := doc.Find("main")
	if main.Length() > 0 {
		return main.Find("p").First().Text()
	}

	return doc.Find("p").First().Text()
}

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var urls []string

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists {
			return
		}

		if strings.HasPrefix(href, "/") {
			hrefStruct, err := url.Parse(href)
			if err != nil {
				return
			}

			absoluteURL := baseURL.ResolveReference(hrefStruct)
			urls = append(urls, absoluteURL.String())
		} else {
			urls = append(urls, href)
		}
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	var images []string

	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if !exists {
			return
		}

		if strings.HasPrefix(src, "/") {
			images = append(images, baseURL.String()+src)
		} else {
			images = append(images, src)
		}
	})

	return images, nil
}

type PageData struct {
	URL            *url.URL
	Heading        string
	FirstParagraph string
	OutgoingLinks  []string
	ImageURLs      []string
}

func extractPageData(html, pageURL string) (*PageData, error) {
	urlStruct, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}

	outgoingLinks, err := getURLsFromHTML(html, urlStruct)
	if err != nil {
		return nil, err
	}

	imageURLs, _ := getImagesFromHTML(html, urlStruct)
	if err != nil {
		return nil, err
	}

	return &PageData{
		URL:            urlStruct,
		Heading:        getHeadingFromHTML(html),
		FirstParagraph: getFirstParagraphFromHTML(html),
		OutgoingLinks:  outgoingLinks,
		ImageURLs:      imageURLs,
	}, nil
}

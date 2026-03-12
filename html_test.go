package main

import (
	"net/url"
	"reflect"
	"testing"
)

func TestGetHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name: "h1",
			html: `
<html>
  <body>
    <h1>Welcome to Boot.dev</h1>
  </body>
</html>
`,
			expected: "Welcome to Boot.dev",
		},
		{
			name: "h2", html: `
<html>
  <body>
    <h2>Welcome to Boot.dev</h2>
  </body>
</html>
`,
			expected: "Welcome to Boot.dev",
		},
		{
			name: "no header",
			html: `
<html>
  <body>
  </body>
</html>
`,
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getHeadingFromHTML(tc.html)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected heading: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name: "main",
			html: `
<html>
  <body>
    <p>Outer</p>
    <main>
      <p>Learn to code by building real projects.</p>
    </main>
  </body>
</html>
`,
			expected: "Learn to code by building real projects.",
		},
		{
			name: "no main",
			html: `
<html>
  <body>
    <p>Learn to code by building real projects.</p>
  </body>
</html>
`,
			expected: "Learn to code by building real projects.",
		},
		{
			name: "No paragraph",
			html: `
<html>
  <body>
  </body>
</html>
`,
			expected: "",
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := getFirstParagraphFromHTML(tc.html)
			if actual != tc.expected {
				t.Errorf("Test %v - %s FAIL: expected heading: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetURLsFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		htmlBody string
		expected []string
	}{
		{
			name:     "absolute",
			inputURL: "https://boot.dev",
			htmlBody: `
<html>
  <body>
		<a href="https://crawler-test.com">Crawler Test</a>
  </body>
</html>
`,
			expected: []string{"https://crawler-test.com"},
		},
		{
			name:     "relative",
			inputURL: "https://boot.dev",
			htmlBody: `
<html>
  <body>
		<a href="/crawler-test">Boot.dev</a>
  </body>
</html>
`,
			expected: []string{"https://boot.dev/crawler-test"},
		},
		{
			name:     "multiple links",
			inputURL: "https://boot.dev",
			htmlBody: `
<html>
  <body>
		<a href="https://crawler-test.com">Crawler Test</a>
		<a href="/crawler-test">Boot.dev</a>
  </body>
</html>
`,
			expected: []string{
				"https://crawler-test.com",
				"https://boot.dev/crawler-test",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - %s FAIL: couldn't parse input URL: %v", i, tc.name, err)
				return
			}

			actual, err := getURLsFromHTML(tc.htmlBody, baseURL)
			if err != nil {
				t.Fatalf("Test %v - %s FAIL: unexpected error: %v", i, tc.name, err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestGetImagesFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		inputURL string
		htmlBody string
		expected []string
	}{
		{
			name:     "absolute",
			inputURL: "https://boot.dev",
			htmlBody: `
<html>
  <body>
		<img src="https://crawler-test.com/logo.png" alt="Logo">
  </body>
</html>
`,
			expected: []string{"https://crawler-test.com/logo.png"},
		},
		{
			name:     "relative",
			inputURL: "https://boot.dev",
			htmlBody: `
<html>
  <body>
		<img src="/logo.png" alt="Logo">
  </body>
</html>
`,
			expected: []string{"https://boot.dev/logo.png"},
		},
		{
			name:     "multiple links",
			inputURL: "https://boot.dev",
			htmlBody: `
<html>
  <body>
		<img src="https://crawler-test.com/logo.png" alt="Logo">
		<img src="/logo.png" alt="Logo">
  </body>
</html>
`,
			expected: []string{
				"https://crawler-test.com/logo.png",
				"https://boot.dev/logo.png",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			baseURL, err := url.Parse(tc.inputURL)
			if err != nil {
				t.Errorf("Test %v - %s FAIL: couldn't parse input URL: %v", i, tc.name, err)
				return
			}

			actual, err := getImagesFromHTML(tc.htmlBody, baseURL)
			if err != nil {
				t.Fatalf("Test %v - %s FAIL: unexpected error: %v", i, tc.name, err)
			}

			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("Test %v - %s FAIL: expected: %v, actual: %v", i, tc.name, tc.expected, actual)
			}
		})
	}
}

func TestExtractPageData(t *testing.T) {
	inputURL := "https://crawler-test.com"
	inputBody := `
<html>
	<body>
    <h1>Test Title</h1>
    <p>This is the first paragraph.</p>
    <a href="/link1">Link 1</a>
    <img src="/image1.jpg" alt="Image 1">
  </body>
</html>`

	urlStruct, _ := url.Parse(inputURL)

	actual, err := extractPageData(inputBody, inputURL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := PageData{
		URL:            urlStruct,
		Heading:        "Test Title",
		FirstParagraph: "This is the first paragraph.",
		OutgoingLinks:  []string{"https://crawler-test.com/link1"},
		ImageURLs:      []string{"https://crawler-test.com/image1.jpg"},
	}

	if !reflect.DeepEqual(*actual, expected) {
		t.Errorf("expected %+v, got %+v", expected, *actual)
	}
}

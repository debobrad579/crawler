# Crawler

A concurrent web crawler written in Go that crawls a website and generates a JSON report with page data.

## Usage

```bash
go run . <url> [max_concurrency] [max_pages]
```

- `url` – The starting URL to crawl (required)
- `max_concurrency` – Max number of concurrent requests (default: 5)
- `max_pages` – Max number of pages to crawl (default: 25)

**Example:**

```bash
go run . https://boot.dev 10 50
```

## Output

Reports are saved as JSON files in the `reports/` directory. Each page entry includes:

- URL
- Heading (h1 or h2)
- First paragraph
- Outgoing links
- Image URLs

## Running Tests

```bash
go test ./...
```

## License

This project is part of the Boot.dev curriculum.

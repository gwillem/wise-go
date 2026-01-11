// Package main implements a documentation crawler that fetches Wise API documentation
// and converts it to Markdown files.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

const (
	baseURL     = "https://docs.wise.com"
	sitemapURL  = "https://docs.wise.com/sitemap.xml"
	outputDir   = "docs/api-reference"
	concurrency = 3
	rateLimit   = 500 * time.Millisecond
)

type Crawler struct {
	client     *http.Client
	visited    map[string]bool
	mu         sync.Mutex
	rateLimiter <-chan time.Time
}

func NewCrawler() *Crawler {
	return &Crawler{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		visited:     make(map[string]bool),
		rateLimiter: time.Tick(rateLimit),
	}
}

func (c *Crawler) markVisited(path string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.visited[path] {
		return false
	}
	c.visited[path] = true
	return true
}

func (c *Crawler) fetch(ctx context.Context, urlStr string) (string, error) {
	<-c.rateLimiter

	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "WiseDocsCrawler/1.0")
	req.Header.Set("Accept", "text/html")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, urlStr)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// extractLinks finds all links to other API reference pages
func extractLinks(htmlContent string, currentURL *url.URL) []string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	var links []string
	var findLinks func(*html.Node)
	findLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href := attr.Val
					// Resolve relative URLs
					parsed, err := url.Parse(href)
					if err != nil {
						continue
					}
					resolved := currentURL.ResolveReference(parsed)

					// Only include links to the API reference section
					if resolved.Host == currentURL.Host && strings.HasPrefix(resolved.Path, "/api-reference") {
						// Normalize the path (remove trailing slash, fragments, query params)
						path := strings.TrimSuffix(resolved.Path, "/")
						if path != "" && !strings.Contains(path, "#") {
							links = append(links, path)
						}
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			findLinks(child)
		}
	}
	findLinks(doc)

	// Deduplicate
	seen := make(map[string]bool)
	var unique []string
	for _, link := range links {
		if !seen[link] {
			seen[link] = true
			unique = append(unique, link)
		}
	}
	return unique
}

// extractMainContent extracts the main documentation content from the HTML
func extractMainContent(htmlContent string) *html.Node {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return nil
	}

	// Look for main content containers (common patterns in documentation sites)
	var mainContent *html.Node
	var findMain func(*html.Node) bool
	findMain = func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			// Check for common main content indicators
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					classes := strings.ToLower(attr.Val)
					if strings.Contains(classes, "main-content") ||
						strings.Contains(classes, "documentation") ||
						strings.Contains(classes, "doc-content") ||
						strings.Contains(classes, "article") ||
						strings.Contains(classes, "content-body") {
						mainContent = n
						return true
					}
				}
				if attr.Key == "role" && attr.Val == "main" {
					mainContent = n
					return true
				}
				if attr.Key == "id" {
					id := strings.ToLower(attr.Val)
					if id == "main" || id == "content" || id == "main-content" {
						mainContent = n
						return true
					}
				}
			}
			if n.Data == "main" || n.Data == "article" {
				mainContent = n
				return true
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if findMain(child) {
				return true
			}
		}
		return false
	}
	findMain(doc)

	if mainContent != nil {
		return mainContent
	}

	// Fallback: return body if no main content found
	var findBody func(*html.Node) *html.Node
	findBody = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "body" {
			return n
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if body := findBody(child); body != nil {
				return body
			}
		}
		return nil
	}
	return findBody(doc)
}

// htmlToMarkdown converts an HTML node tree to Markdown
func htmlToMarkdown(n *html.Node) string {
	if n == nil {
		return ""
	}

	var sb strings.Builder
	convertNode(&sb, n, 0, false)

	return cleanupMarkdown(sb.String())
}

// cleanupMarkdown applies post-processing transformations to clean up markdown
func cleanupMarkdown(result string) string {
	// Clean up excessive newlines
	result = regexp.MustCompile(`\n{3,}`).ReplaceAllString(result, "\n\n")

	// Remove empty anchor links from headings: [](#anchor-name) -> nothing
	result = regexp.MustCompile(`\[]\(#[^)]*\)`).ReplaceAllString(result, "")

	// Remove "Copy" button text from end of headings
	result = regexp.MustCompile(`(?m)^(#{1,6} .+?)Copy\s*$`).ReplaceAllString(result, "$1")

	// Fix code fence closing on same line as content: }``` -> }\n```
	result = regexp.MustCompile("([^\n])```").ReplaceAllString(result, "$1\n```")

	// Fix code fence with content immediately after: ```\n\n## -> ```\n\n##
	result = regexp.MustCompile("```([A-Z#])").ReplaceAllString(result, "```\n\n$1")

	// Fix spaces inside bold markers: ** text** -> **text**
	boldSpaceRe := regexp.MustCompile(`\*\* ([^*\n]+)\*\*`)
	for boldSpaceRe.MatchString(result) {
		result = boldSpaceRe.ReplaceAllString(result, "**$1**")
	}

	// Fix spaces inside inline code: ` code` -> `code`
	// Run multiple times to catch all occurrences
	for i := 0; i < 10; i++ {
		result = regexp.MustCompile("` ([^`\n]+)`").ReplaceAllString(result, "`$1`")
	}

	// Fix missing spaces after bold text (e.g., **text**word -> **text** word)
	// Use [^*\n]+ to prevent matching across newlines (greedy match would capture too much)
	result = regexp.MustCompile(`\*\*([^*\n]+)\*\*([a-zA-Z])`).ReplaceAllString(result, "**$1** $2")

	// Fix missing spaces after inline code (e.g., `code`word -> `code` word)
	// Use [^`\s,]+ to avoid matching across multiple inline code spans
	result = regexp.MustCompile("`([^`\\s,]+)`([a-zA-Z])").ReplaceAllString(result, "`$1` $2")

	// Convert "Fields**field**" pattern to proper formatting with newline
	result = regexp.MustCompile(`(?m)^Fields\*\*`).ReplaceAllString(result, "**")

	// Convert "Request**field**" and "Response**field**" patterns
	result = regexp.MustCompile(`(?m)^Request\*\*`).ReplaceAllString(result, "### Request\n\n**")
	result = regexp.MustCompile(`(?m)^Response\*\*`).ReplaceAllString(result, "### Response\n\n**")
	result = regexp.MustCompile(`(?m)^Request Parameters\*\*`).ReplaceAllString(result, "### Request Parameters\n\n**")

	// Remove "Was this helpful" section and anything after
	if idx := strings.Index(result, "Was this helpful?"); idx != -1 {
		result = result[:idx]
	}
	if idx := strings.Index(result, "#### Was this helpful"); idx != -1 {
		result = result[:idx]
	}

	// Remove trailing incomplete headers (e.g., #### at end of file)
	result = regexp.MustCompile(`(?m)\n#{1,6}\s*$`).ReplaceAllString(result, "")

	// Break apart squished operation lists: [METHOD/path](link)Description[METHOD -> newlines
	// Pattern: ](link)Text[  ->  ](link)\n\nText\n\n[
	result = regexp.MustCompile(`\]\(([^)]+)\)([A-Z][^[]+)\[`).ReplaceAllString(result, "]($1)\n\n$2\n\n[")

	// Clean up any quadruple+ newlines created
	result = regexp.MustCompile(`\n{3,}`).ReplaceAllString(result, "\n\n")

	result = strings.TrimSpace(result)
	return result
}

func convertNode(sb *strings.Builder, n *html.Node, depth int, inPre bool) {
	if n == nil {
		return
	}

	switch n.Type {
	case html.TextNode:
		text := n.Data
		if !inPre {
			// Collapse whitespace outside of pre blocks
			text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
		}
		sb.WriteString(text)
		return

	case html.ElementNode:
		// Skip elements we don't want in documentation
		if shouldSkipElement(n) {
			return
		}

		switch n.Data {
		case "h1":
			sb.WriteString("\n\n# ")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("\n\n")
			return
		case "h2":
			sb.WriteString("\n\n## ")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("\n\n")
			return
		case "h3":
			sb.WriteString("\n\n### ")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("\n\n")
			return
		case "h4":
			sb.WriteString("\n\n#### ")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("\n\n")
			return
		case "h5":
			sb.WriteString("\n\n##### ")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("\n\n")
			return
		case "h6":
			sb.WriteString("\n\n###### ")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("\n\n")
			return
		case "p":
			sb.WriteString("\n\n")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("\n\n")
			return
		case "br":
			sb.WriteString("\n")
			return
		case "hr":
			sb.WriteString("\n\n---\n\n")
			return
		case "strong", "b":
			sb.WriteString("**")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("**")
			return
		case "em", "i":
			sb.WriteString("*")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("*")
			return
		case "code":
			if !inPre {
				sb.WriteString("`")
				convertChildren(sb, n, depth, true)
				sb.WriteString("`")
			} else {
				convertChildren(sb, n, depth, true)
			}
			return
		case "pre":
			sb.WriteString("\n\n```")
			// Try to get language from class
			lang := getCodeLanguage(n)
			sb.WriteString(lang)
			sb.WriteString("\n")
			convertChildren(sb, n, depth, true)
			sb.WriteString("\n```\n\n")
			return
		case "a":
			href := getAttr(n, "href")
			sb.WriteString("[")
			convertChildren(sb, n, depth, inPre)
			sb.WriteString("](")
			sb.WriteString(href)
			sb.WriteString(")")
			return
		case "img":
			alt := getAttr(n, "alt")
			src := getAttr(n, "src")
			sb.WriteString("![")
			sb.WriteString(alt)
			sb.WriteString("](")
			sb.WriteString(src)
			sb.WriteString(")")
			return
		case "ul":
			sb.WriteString("\n")
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.ElementNode && child.Data == "li" {
					sb.WriteString("\n- ")
					convertChildren(sb, child, depth+1, inPre)
				}
			}
			sb.WriteString("\n")
			return
		case "ol":
			sb.WriteString("\n")
			i := 1
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.ElementNode && child.Data == "li" {
					fmt.Fprintf(sb, "\n%d. ", i)
					convertChildren(sb, child, depth+1, inPre)
					i++
				}
			}
			sb.WriteString("\n")
			return
		case "blockquote":
			sb.WriteString("\n\n> ")
			text := getTextContent(n)
			text = strings.ReplaceAll(text, "\n", "\n> ")
			sb.WriteString(text)
			sb.WriteString("\n\n")
			return
		case "table":
			convertTable(sb, n)
			return
		case "div", "section", "article", "main", "span":
			convertChildren(sb, n, depth, inPre)
			return
		case "script", "style", "nav", "header", "footer", "aside", "noscript":
			// Skip these entirely
			return
		default:
			convertChildren(sb, n, depth, inPre)
		}
	default:
		convertChildren(sb, n, depth, inPre)
	}
}

func shouldSkipElement(n *html.Node) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.ToLower(attr.Val)
			// Skip navigation, sidebars, footers, etc.
			if strings.Contains(classes, "nav") ||
				strings.Contains(classes, "sidebar") ||
				strings.Contains(classes, "footer") ||
				strings.Contains(classes, "header") ||
				strings.Contains(classes, "menu") ||
				strings.Contains(classes, "breadcrumb") ||
				strings.Contains(classes, "toc") ||
				strings.Contains(classes, "pagination") ||
				strings.Contains(classes, "cookie") ||
				strings.Contains(classes, "banner") {
				return true
			}
		}
		if attr.Key == "role" {
			role := strings.ToLower(attr.Val)
			if role == "navigation" || role == "banner" || role == "contentinfo" {
				return true
			}
		}
		if attr.Key == "hidden" {
			return true
		}
	}
	return false
}

func convertChildren(sb *strings.Builder, n *html.Node, depth int, inPre bool) {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		convertNode(sb, child, depth, inPre)
	}
}

func convertTable(sb *strings.Builder, table *html.Node) {
	var headers []string
	var rows [][]string

	var processTable func(*html.Node)
	processTable = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "thead":
				for child := n.FirstChild; child != nil; child = child.NextSibling {
					if child.Type == html.ElementNode && child.Data == "tr" {
						for th := child.FirstChild; th != nil; th = th.NextSibling {
							if th.Type == html.ElementNode && (th.Data == "th" || th.Data == "td") {
								headers = append(headers, strings.TrimSpace(getTextContent(th)))
							}
						}
					}
				}
			case "tbody":
				for child := n.FirstChild; child != nil; child = child.NextSibling {
					if child.Type == html.ElementNode && child.Data == "tr" {
						var row []string
						for td := child.FirstChild; td != nil; td = td.NextSibling {
							if td.Type == html.ElementNode && (td.Data == "td" || td.Data == "th") {
								row = append(row, strings.TrimSpace(getTextContent(td)))
							}
						}
						if len(row) > 0 {
							rows = append(rows, row)
						}
					}
				}
			case "tr":
				// Handle tables without thead/tbody
				if n.Parent != nil && n.Parent.Data == "table" {
					var row []string
					isHeader := false
					for td := n.FirstChild; td != nil; td = td.NextSibling {
						if td.Type == html.ElementNode {
							if td.Data == "th" {
								isHeader = true
							}
							if td.Data == "td" || td.Data == "th" {
								row = append(row, strings.TrimSpace(getTextContent(td)))
							}
						}
					}
					if len(row) > 0 {
						if isHeader && len(headers) == 0 {
							headers = row
						} else {
							rows = append(rows, row)
						}
					}
				}
			}
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			processTable(child)
		}
	}
	processTable(table)

	if len(headers) == 0 && len(rows) == 0 {
		return
	}

	sb.WriteString("\n\n")

	// Write headers
	if len(headers) > 0 {
		sb.WriteString("| ")
		sb.WriteString(strings.Join(headers, " | "))
		sb.WriteString(" |\n")
		sb.WriteString("|")
		for range headers {
			sb.WriteString(" --- |")
		}
		sb.WriteString("\n")
	}

	// Write rows
	for _, row := range rows {
		sb.WriteString("| ")
		sb.WriteString(strings.Join(row, " | "))
		sb.WriteString(" |\n")
	}
	sb.WriteString("\n")
}

func getAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getCodeLanguage(n *html.Node) string {
	// Check class attribute on pre element
	class := getAttr(n, "class")

	// Check data-language attribute
	if lang := getAttr(n, "data-language"); lang != "" {
		return strings.ToLower(lang)
	}
	if lang := getAttr(n, "data-lang"); lang != "" {
		return strings.ToLower(lang)
	}

	// Check child code element for class or data attributes
	var codeNode *html.Node
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.Data == "code" {
			codeNode = child
			if class == "" {
				class = getAttr(child, "class")
			}
			if lang := getAttr(child, "data-language"); lang != "" {
				return strings.ToLower(lang)
			}
			if lang := getAttr(child, "data-lang"); lang != "" {
				return strings.ToLower(lang)
			}
			break
		}
	}

	// Common patterns: language-json, lang-json, hljs json, etc.
	patterns := []string{
		`language-(\w+)`,
		`lang-(\w+)`,
		`\bhljs\s+(\w+)`,
		`\b(json|javascript|typescript|python|ruby|go|bash|shell|curl|xml|html|css|yaml|sql)\b`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(strings.ToLower(class)); len(matches) > 1 {
			return matches[1]
		}
	}

	// Infer language from content
	content := strings.TrimSpace(getTextContent(n))
	if codeNode != nil {
		content = strings.TrimSpace(getTextContent(codeNode))
	}

	return inferLanguageFromContent(content)
}

func inferLanguageFromContent(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return ""
	}

	// Check for curl commands
	if strings.HasPrefix(content, "curl ") || strings.HasPrefix(content, "curl\n") {
		return "bash"
	}

	// Check for shell commands
	if strings.HasPrefix(content, "$ ") || strings.HasPrefix(content, "# ") {
		return "bash"
	}

	// Check for JSON (starts with { or [)
	if (strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}")) ||
		(strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]")) {
		return "json"
	}

	// Check for HTTP methods (likely HTTP request examples)
	httpMethods := []string{"GET ", "POST ", "PUT ", "PATCH ", "DELETE ", "HEAD ", "OPTIONS "}
	for _, method := range httpMethods {
		if strings.HasPrefix(content, method) {
			return "http"
		}
	}

	// Check for XML
	if strings.HasPrefix(content, "<?xml") || strings.HasPrefix(content, "<") {
		return "xml"
	}

	return ""
}

func getTextContent(n *html.Node) string {
	var sb strings.Builder
	var extract func(*html.Node)
	extract = func(node *html.Node) {
		if node.Type == html.TextNode {
			sb.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			extract(child)
		}
	}
	extract(n)
	return sb.String()
}

// fetchSitemap fetches the sitemap.xml and extracts URLs matching the prefix
func (c *Crawler) fetchSitemap(ctx context.Context, prefix string) ([]string, error) {
	<-c.rateLimiter

	req, err := http.NewRequestWithContext(ctx, "GET", sitemapURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "WiseDocsCrawler/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("sitemap HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse sitemap XML - simple regex approach for <loc> tags
	locPattern := regexp.MustCompile(`<loc>([^<]+)</loc>`)
	matches := locPattern.FindAllStringSubmatch(string(body), -1)

	var urls []string
	for _, match := range matches {
		if len(match) > 1 {
			urlStr := match[1]
			if strings.Contains(urlStr, prefix) {
				urls = append(urls, urlStr)
			}
		}
	}

	return urls, nil
}

func pathToFilename(urlPath string) string {
	// Remove /api-reference prefix and convert to filename
	path := strings.TrimPrefix(urlPath, "/api-reference")
	path = strings.TrimPrefix(path, "/")

	if path == "" {
		return "index.md"
	}

	// Replace slashes with dashes or create directory structure
	path = strings.ReplaceAll(path, "/", "-")
	// Clean up any special characters
	path = regexp.MustCompile(`[^a-zA-Z0-9\-_]`).ReplaceAllString(path, "-")
	path = regexp.MustCompile(`-+`).ReplaceAllString(path, "-")
	path = strings.Trim(path, "-")

	return path + ".md"
}

func (c *Crawler) crawl(ctx context.Context) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Fetch URLs from sitemap
	log.Printf("Fetching sitemap...")
	urls, err := c.fetchSitemap(ctx, "/api-reference")
	if err != nil {
		return fmt.Errorf("failed to fetch sitemap: %w", err)
	}
	log.Printf("Found %d API reference pages in sitemap", len(urls))

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrency)

	for _, fullURL := range urls {
		parsed, err := url.Parse(fullURL)
		if err != nil {
			log.Printf("Invalid URL %s: %v", fullURL, err)
			continue
		}

		if !c.markVisited(parsed.Path) {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(urlStr string, urlPath string) {
			defer wg.Done()
			defer func() { <-sem }()

			log.Printf("Fetching: %s", urlStr)

			htmlContent, err := c.fetch(ctx, urlStr)
			if err != nil {
				log.Printf("Error fetching %s: %v", urlStr, err)
				return
			}

			// Extract main content and convert to markdown
			mainContent := extractMainContent(htmlContent)
			markdown := htmlToMarkdown(mainContent)

			if strings.TrimSpace(markdown) == "" {
				log.Printf("No content extracted from %s", urlStr)
				return
			}

			// Add source URL as header
			markdown = fmt.Sprintf("<!-- Source: %s -->\n\n%s", urlStr, markdown)

			// Write to file
			filename := pathToFilename(urlPath)
			outPath := filepath.Join(outputDir, filename)

			if err := os.WriteFile(outPath, []byte(markdown), 0644); err != nil {
				log.Printf("Error writing %s: %v", outPath, err)
				return
			}

			log.Printf("Wrote: %s", outPath)
		}(fullURL, parsed.Path)
	}

	wg.Wait()
	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	crawler := NewCrawler()

	log.Printf("Starting crawl using sitemap: %s", sitemapURL)
	log.Printf("Output directory: %s", outputDir)

	if err := crawler.crawl(ctx); err != nil {
		log.Fatalf("Crawl failed: %v", err)
	}

	log.Printf("Crawl completed. Visited %d pages.", len(crawler.visited))
}

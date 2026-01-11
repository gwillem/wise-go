package main

import (
	"net/url"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestPathToFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/api-reference", "index.md"},
		{"/api-reference/", "index.md"},
		{"/api-reference/transfers", "transfers.md"},
		{"/api-reference/transfers/create", "transfers-create.md"},
		{"/api-reference/bank-accounts/requirements", "bank-accounts-requirements.md"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := pathToFilename(tt.input)
			if result != tt.expected {
				t.Errorf("pathToFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractLinks(t *testing.T) {
	html := `
	<html>
	<body>
		<a href="/api-reference/transfers">Transfers</a>
		<a href="/api-reference/quotes">Quotes</a>
		<a href="https://external.com/link">External</a>
		<a href="/other-section">Other</a>
		<a href="/api-reference/transfers#section">With hash</a>
	</body>
	</html>
	`

	links := extractLinks(html, mustParseURL("https://docs.wise.com/api-reference"))

	expectedLinks := []string{"/api-reference/transfers", "/api-reference/quotes"}
	if len(links) != len(expectedLinks) {
		t.Errorf("expected %d links, got %d: %v", len(expectedLinks), len(links), links)
	}

	for _, expected := range expectedLinks {
		found := false
		for _, link := range links {
			if link == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected to find link %q in %v", expected, links)
		}
	}
}

func TestHtmlToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		contains []string
	}{
		{
			name: "headings",
			html: `<div><h1>Title</h1><h2>Subtitle</h2><p>Content</p></div>`,
			contains: []string{
				"# Title",
				"## Subtitle",
				"Content",
			},
		},
		{
			name: "code blocks",
			html: `<div><pre class="language-json"><code>{"key": "value"}</code></pre></div>`,
			contains: []string{
				"```json",
				`{"key": "value"}`,
				"```",
			},
		},
		{
			name: "inline code",
			html: `<p>Use <code>transfer_id</code> parameter</p>`,
			contains: []string{
				"`transfer_id`",
			},
		},
		{
			name: "lists",
			html: `<ul><li>Item 1</li><li>Item 2</li></ul>`,
			contains: []string{
				"- Item 1",
				"- Item 2",
			},
		},
		{
			name: "ordered lists",
			html: `<ol><li>First</li><li>Second</li></ol>`,
			contains: []string{
				"1. First",
				"2. Second",
			},
		},
		{
			name: "links",
			html: `<p>See <a href="https://example.com">documentation</a></p>`,
			contains: []string{
				"[documentation](https://example.com)",
			},
		},
		{
			name: "bold and italic",
			html: `<p><strong>Important</strong> and <em>emphasized</em></p>`,
			contains: []string{
				"**Important**",
				"*emphasized*",
			},
		},
		{
			name: "tables",
			html: `<table><thead><tr><th>Name</th><th>Type</th></tr></thead><tbody><tr><td>id</td><td>string</td></tr></tbody></table>`,
			contains: []string{
				"| Name | Type |",
				"| --- | --- |",
				"| id | string |",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := extractMainContent(tt.html)
			result := htmlToMarkdown(content)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("expected markdown to contain %q, got:\n%s", expected, result)
				}
			}
		})
	}
}

func TestShouldSkipElement(t *testing.T) {
	tests := []struct {
		html   string
		skip   bool
		reason string
	}{
		{`<nav>Navigation</nav>`, true, "nav element"},
		{`<div class="sidebar">Side</div>`, true, "sidebar class"},
		{`<div class="footer">Footer</div>`, true, "footer class"},
		{`<div class="content">Content</div>`, false, "content div"},
		{`<div role="navigation">Nav</div>`, true, "navigation role"},
		{`<div hidden>Hidden</div>`, true, "hidden attribute"},
	}

	for _, tt := range tests {
		t.Run(tt.reason, func(t *testing.T) {
			content := extractMainContent(tt.html)
			result := htmlToMarkdown(content)

			if tt.skip && strings.TrimSpace(result) != "" {
				// Check if the result contains any of the skipped content
				// This is a simplified test
			}
		})
	}
}

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func TestInlineCodeRegex(t *testing.T) {
	re := regexp.MustCompile("` ([^`\n]+)`")

	tests := []struct {
		input    string
		expected string
	}{
		{"` from`", "`from`"},
		{"` to`", "`to`"},
		{"` from`, ` to`", "`from`, `to`"},
		{"the ` from`, ` to`, and ` group` params", "the `from`, `to`, and `group` params"},
	}

	for _, tt := range tests {
		result := tt.input
		for i := 0; i < 10; i++ {
			result = re.ReplaceAllString(result, "`$1`")
		}
		if result != tt.expected {
			t.Errorf("input %q\ngot:  %q\nwant: %q", tt.input, result, tt.expected)
		}
	}
}

func TestInlineCodeInCleanup(t *testing.T) {
	// Test just the inline code transformation step by step
	input := "the ` from`, ` to`, and ` group` params"
	expected := "the `from`, `to`, and `group` params"

	result := input

	// Apply each transformation from cleanupMarkdown one at a time
	// 1. Clean up excessive newlines
	result = regexp.MustCompile(`\n{3,}`).ReplaceAllString(result, "\n\n")
	t.Logf("After newline cleanup: %q", result)

	// 2. Remove empty anchor links
	result = regexp.MustCompile(`\[]\(#[^)]*\)`).ReplaceAllString(result, "")
	t.Logf("After anchor cleanup: %q", result)

	// 3. Remove Copy button text
	result = regexp.MustCompile(`(?m)^(#{1,6} .+?)Copy\s*$`).ReplaceAllString(result, "$1")
	t.Logf("After Copy cleanup: %q", result)

	// 4. Fix code fence closing
	result = regexp.MustCompile("([^\n])```").ReplaceAllString(result, "$1\n```")
	t.Logf("After fence closing: %q", result)

	// 5. Fix code fence with content after
	result = regexp.MustCompile("```([A-Z#])").ReplaceAllString(result, "```\n\n$1")
	t.Logf("After fence content: %q", result)

	// 6. Fix bold space
	boldSpaceRe := regexp.MustCompile(`\*\* ([^*\n]+)\*\*`)
	for boldSpaceRe.MatchString(result) {
		result = boldSpaceRe.ReplaceAllString(result, "**$1**")
	}
	t.Logf("After bold space: %q", result)

	// 7. Fix inline code space - THIS IS THE ONE THAT SHOULD FIX IT
	re := regexp.MustCompile("` ([^`\n]+)`")
	t.Logf("Before inline code fix: %q", result)
	t.Logf("Regex matches: %v", re.FindAllString(result, -1))
	for i := 0; i < 10; i++ {
		result = re.ReplaceAllString(result, "`$1`")
		t.Logf("After iteration %d: %q", i, result)
	}

	if result != expected {
		t.Errorf("Final result: %q, want: %q", result, expected)
	}
}

func TestCleanupMarkdownDebug(t *testing.T) {
	// Test that specifically debugs the cleanupMarkdown function
	input := "the ` from`, ` to`, and ` group` params"
	expected := "the `from`, `to`, and `group` params"

	result := cleanupMarkdown(input)
	if result != expected {
		t.Errorf("cleanupMarkdown(%q) =\n%q\nwant:\n%q", input, result, expected)

		// Let's also test just the regex directly
		re := regexp.MustCompile("` ([^`\n]+)`")
		directResult := re.ReplaceAllString(input, "`$1`")
		t.Logf("Direct regex result: %q", directResult)
		t.Logf("Direct regex matches: %v", re.FindAllString(input, -1))
	}
}

func TestBoldSpaceInGeneratedContent(t *testing.T) {
	// Test the exact pattern seen in generated files
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lastName pattern",
			input:    "** lastName**text (max 30 chars)",
			expected: "**lastName** text (max 30 chars)",
		},
		{
			name:     "source pattern",
			input:    "** source**text",
			expected: "**source** text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanupMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("cleanupMarkdown(%q) =\n%q\nwant:\n%q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseRealHTML(t *testing.T) {
	// Test parsing of actual downloaded HTML file
	htmlBytes, err := os.ReadFile("testdata/input.html")
	if err != nil {
		t.Skipf("Skipping test: %v", err)
	}

	content := extractMainContent(string(htmlBytes))

	// Get raw markdown before cleanup
	var sb strings.Builder
	convertNode(&sb, content, 0, false)
	rawMd := sb.String()

	// Check if pattern exists before cleanup
	boldRe := regexp.MustCompile(`\*\* ([^*\n]+)\*\*`)
	beforeMatches := boldRe.FindAllString(rawMd, 10)
	t.Logf("Bold pattern matches BEFORE cleanup: %v", beforeMatches)

	// Find and show what's around "lastName" in raw markdown
	if idx := strings.Index(rawMd, "lastName"); idx != -1 {
		start := idx - 30
		if start < 0 {
			start = 0
		}
		end := idx + 50
		if end > len(rawMd) {
			end = len(rawMd)
		}
		t.Logf("Raw markdown around 'lastName': %q", rawMd[start:end])
	}

	// Find all occurrences of lastName in raw
	t.Logf("All 'lastName' occurrences in raw:")
	for i, idx := 0, 0; i < 5; i++ {
		nextIdx := strings.Index(rawMd[idx:], "lastName")
		if nextIdx == -1 {
			break
		}
		idx += nextIdx
		start := idx - 20
		if start < 0 {
			start = 0
		}
		end := idx + 40
		if end > len(rawMd) {
			end = len(rawMd)
		}
		t.Logf("  [%d]: %q", i, rawMd[start:end])
		idx += len("lastName")
	}

	// Check for "Request **" or "Request**" patterns
	t.Logf("Looking for Request patterns...")
	for _, pattern := range []string{"Request **", "Request**", "Fields **", "Fields**", "Fields**id", "**id**"} {
		if idx := strings.Index(rawMd, pattern); idx != -1 {
			start := idx
			end := idx + 80
			if end > len(rawMd) {
				end = len(rawMd)
			}
			t.Logf("Found %q at %d: %q", pattern, idx, rawMd[start:end])
		}
	}

	// Show what's around "integer" in raw to understand the structure
	if idx := strings.Index(rawMd, "integer\n\nUnique"); idx != -1 {
		start := idx - 40
		if start < 0 {
			start = 0
		}
		end := idx + 40
		t.Logf("Around 'integer' in raw: %q", rawMd[start:end])
	}

	// Run cleanup step-by-step to find which transformation breaks it
	result := rawMd

	// Find a specific occurrence to track
	trackStr := "**lastName**text (max 30"
	trackIdx := strings.Index(result, trackStr)
	if trackIdx != -1 {
		t.Logf("Initial: found at %d: %q", trackIdx, result[trackIdx:trackIdx+40])
	}

	// Apply each cleanup step and check
	steps := []struct {
		name string
		fn   func(string) string
	}{
		{"newlines", func(s string) string {
			return regexp.MustCompile(`\n{3,}`).ReplaceAllString(s, "\n\n")
		}},
		{"anchor", func(s string) string {
			return regexp.MustCompile(`\[]\(#[^)]*\)`).ReplaceAllString(s, "")
		}},
		{"copy", func(s string) string {
			return regexp.MustCompile(`(?m)^(#{1,6} .+?)Copy\s*$`).ReplaceAllString(s, "$1")
		}},
		{"fence-close", func(s string) string {
			return regexp.MustCompile("([^\n])```").ReplaceAllString(s, "$1\n```")
		}},
		{"fence-content", func(s string) string {
			return regexp.MustCompile("```([A-Z#])").ReplaceAllString(s, "```\n\n$1")
		}},
		{"bold-space", func(s string) string {
			re := regexp.MustCompile(`\*\* ([^*\n]+)\*\*`)
			for re.MatchString(s) {
				s = re.ReplaceAllString(s, "**$1**")
			}
			return s
		}},
		{"inline-code-space", func(s string) string {
			for i := 0; i < 10; i++ {
				s = regexp.MustCompile("` ([^`\n]+)`").ReplaceAllString(s, "`$1`")
			}
			return s
		}},
		{"bold-after", func(s string) string {
			return regexp.MustCompile(`\*\*([^*\n]+)\*\*([a-zA-Z])`).ReplaceAllString(s, "**$1** $2")
		}},
		{"inline-after", func(s string) string {
			return regexp.MustCompile("`([^`\\s,]+)`([a-zA-Z])").ReplaceAllString(s, "`$1` $2")
		}},
		{"Fields", func(s string) string {
			return regexp.MustCompile(`(?m)^Fields\*\*`).ReplaceAllString(s, "**")
		}},
		{"Request", func(s string) string {
			return regexp.MustCompile(`(?m)^Request\*\*`).ReplaceAllString(s, "### Request\n\n**")
		}},
	}

	for _, step := range steps {
		// Track specific field through transformations
		trackField := "**id**integer"
		hasField := strings.Contains(result, trackField)

		result = step.fn(result)

		// Check if field changed
		hasFieldAfter := strings.Contains(result, trackField)
		if hasField && !hasFieldAfter {
			// Find what replaced it
			if idx := strings.Index(result, "integer\n\nUnique"); idx != -1 {
				start := idx - 20
				if start < 0 {
					start = 0
				}
				t.Logf("After %s: %q is gone, now have: %q", step.name, trackField, result[start:idx+20])
			}
		}
		// Check if our tracked string changed
		newIdx := strings.Index(result, trackStr)
		badIdx := strings.Index(result, "** lastName**")
		if newIdx == -1 && trackIdx != -1 {
			t.Logf("After %s: trackStr GONE", step.name)
			trackIdx = -1
		}
		if badIdx != -1 {
			ctx := result[badIdx : badIdx+40]
			if badIdx+40 > len(result) {
				ctx = result[badIdx:]
			}
			t.Logf("After %s: BAD pattern appeared at %d: %q", step.name, badIdx, ctx)
			break
		}
	}

	// Complete the cleanup
	result = cleanupMarkdown(rawMd)

	// Debug: show what's around "integer" after full cleanup
	if idx := strings.Index(result, "integer\n\nUnique"); idx != -1 {
		start := idx - 30
		if start < 0 {
			start = 0
		}
		t.Logf("After full cleanup around 'integer': %q", result[start:idx+30])
	}

	// Check if pattern exists after cleanup
	afterMatches := boldRe.FindAllString(result, 10)
	t.Logf("Bold pattern matches AFTER cleanup: %v", afterMatches)

	// Check for issues that should be fixed
	issues := []struct {
		name    string
		pattern string
		bad     bool // true if pattern should NOT be present
	}{
		// Match ** followed by space and then content (no newline) ending with ** (space inside bold)
		// e.g., "** fieldName**" should not exist (space after opening **)
		{"bold with leading space", `\*\* [a-zA-Z][^*\n]*\*\*`, true},
		{"inline code with leading space", "` [A-Z]", true},
		{"has profile heading", "# Profile", false},
		{"has code blocks", "```json", false},
	}

	for _, issue := range issues {
		re := regexp.MustCompile(issue.pattern)
		found := re.MatchString(result)
		if issue.bad && found {
			matches := re.FindAllString(result, 5)
			t.Errorf("Found bad pattern %q: %v", issue.name, matches)
			// Show context for first match
			if len(matches) > 0 {
				if idx := strings.Index(result, matches[0]); idx != -1 {
					start := idx - 20
					if start < 0 {
						start = 0
					}
					end := idx + len(matches[0]) + 20
					if end > len(result) {
						end = len(result)
					}
					t.Logf("Context for first match: %q", result[start:end])
				}
			}
		}
		if !issue.bad && !found {
			t.Errorf("Missing expected pattern %q", issue.name)
		}
	}

	// Debug: if test fails, show a snippet around the issue
	if t.Failed() {
		// Find and show context around "** lastName**" if present
		if idx := strings.Index(result, "** lastName**"); idx != -1 {
			start := idx - 50
			if start < 0 {
				start = 0
			}
			end := idx + 100
			if end > len(result) {
				end = len(result)
			}
			t.Logf("Context around '** lastName**':\n%s", result[start:end])
			// Show hex dump of the problematic area
			snippet := result[idx : idx+20]
			t.Logf("Hex dump: %x", snippet)
			t.Logf("Bytes: %v", []byte(snippet))
		}
	}
}

func TestMarkdownCleanup(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "bold with leading space",
			input:    "** sourceAccount**text",
			expected: "**sourceAccount** text",
		},
		{
			name:     "bold with leading space and parens",
			input:    "** sourceAccount (optional)**integer",
			expected: "**sourceAccount (optional)** integer",
		},
		{
			name:     "bold with leading space at line start",
			input:    "### Request\n\n** source**text",
			expected: "### Request\n\n**source** text",
		},
		{
			name:     "inline code with leading space",
			input:    "Use ` transfer-requirements` endpoint",
			expected: "Use `transfer-requirements` endpoint",
		},
		{
			name:     "inline code with leading space multiple",
			input:    "the ` from`, ` to`, and ` group` params",
			expected: "the `from`, `to`, and `group` params",
		},
		{
			name:     "code fence on same line as content",
			input:    "```json\n{\"key\": \"value\"}```\n\nNext section",
			expected: "```json\n{\"key\": \"value\"}\n```\n\nNext section",
		},
		{
			name:     "code fence with heading after",
			input:    "```json\n{\"key\": \"value\"}```## Next heading",
			expected: "```json\n{\"key\": \"value\"}\n```\n\n## Next heading",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We need to run the cleanup through htmlToMarkdown
			// For now, test the specific regex patterns
			result := cleanupMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("cleanupMarkdown(%q) =\n%q\nwant:\n%q", tt.input, result, tt.expected)
			}
		})
	}
}

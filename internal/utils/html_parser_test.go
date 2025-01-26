package utils

import (
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestExtractTitle(t *testing.T) {
	htmlContent := `<html><head><title>Test Page</title></head><body></body></html>`
	expectedTitle := "Test Page"

	title := ExtractTitle(htmlContent)
	assert.Equal(t, expectedTitle, title, "Extracted title should match expected value")
}

func TestExtractTitle_Empty(t *testing.T) {
	htmlContent := `<html><head></head><body></body></html>`

	title := ExtractTitle(htmlContent)
	assert.Equal(t, "", title, "Extracted title should be empty when no title tag exists")
}

func TestDetectHTMLVersion(t *testing.T) {
	testCases := []struct {
		htmlContent string
		expected    string
	}{
		{"<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\">", "HTML 4.01"},
		{"<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 3.2 Final//EN\">", "HTML 3.2"},
		{"<!DOCTYPE html>", "HTML5"},
		{"<html><head></head><body></body></html>", "Unknown HTML version"},
	}

	for _, testCase := range testCases {
		version := DetectHTMLVersion(testCase.htmlContent)
		assert.Equal(t, testCase.expected, version, "Detected HTML version should match expected value")
	}
}

func TestCountHeadings(t *testing.T) {
	htmlContent := `<html><body>
		<h1>Heading 1</h1>
		<h2>Heading 2</h2>
		<h2>Heading 2 again</h2>
		<h3>Heading 3</h3>
	</body></html>`

	expected := map[string]int{
		"h1": 1,
		"h2": 2,
		"h3": 1,
	}

	headings := CountHeadings(htmlContent)
	assert.Equal(t, expected, headings, "Headings count should match expected values")
}

func TestCountLinksConcurrently(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Mock HTTP responses
	httpmock.RegisterResponder("HEAD", "http://example.com",
		httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("HEAD", "http://internal-site.com",
		httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("HEAD", "/internal",
		httpmock.NewStringResponder(200, ""))
	httpmock.RegisterResponder("HEAD", "https://broken-link-test.com/",
		httpmock.NewStringResponder(404, ""))

	htmlContent := `<html>
		<body>
			<a href="http://example.com">External Link</a>
			<a href="http://internal-site.com">Internal Link</a>
			<a href="/internal">Internal Link 2</a>
			<a href="https://broken-link-test.com/">Broken Link</a>
		</body>
	</html>`

	internal, external, inaccessible, err := CountLinksConcurrently("http://internal-site.com", htmlContent)

	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, 2, internal, "Should count 1 internal link")
	assert.Equal(t, 2, external, "Should count 1 external link")
	assert.Equal(t, 1, inaccessible, "Should count 1 inaccessible link")
}

func TestContainsLoginForm(t *testing.T) {
	htmlWithLogin := `<html><body><form><input type="password"></form></body></html>`
	htmlWithoutLogin := `<html><body><form><input type="text"></form></body></html>`

	assert.True(t, ContainsLoginForm(htmlWithLogin), "Should detect login form")
	assert.False(t, ContainsLoginForm(htmlWithoutLogin), "Should not detect login form")
}
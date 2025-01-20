package utils

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// ExtractTitle retrieves the title from the given HTML content
func ExtractTitle(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return ""
	}

	var title string
	var parse func(*html.Node)
	parse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" && n.FirstChild != nil {
			title = n.FirstChild.Data
		}
		for c := n.FirstChild; c != nil && title == ""; c = c.NextSibling {
			parse(c)
		}
	}
	parse(doc)

	return strings.TrimSpace(title)
}

func DetectHTMLVersion(htmlContent string) string {
	doctypeRegex := regexp.MustCompile(`(?i)<!DOCTYPE\s+([^>]+)>`)
	matches := doctypeRegex.FindStringSubmatch(htmlContent)

	if len(matches) < 2 {
		return "Unknown HTML version"
	}

	doctype := strings.ToLower(strings.TrimSpace(matches[1]))

	switch {
	case strings.Contains(doctype, "html public \"-//w3c//dtd html 2.0//en\""):
		return "HTML 2.0"
	case strings.Contains(doctype, "html public \"-//w3c//dtd html 3.2 final//en\""):
		return "HTML 3.2"
	case strings.Contains(doctype, "html public \"-//w3c//dtd html 4.01//en\""):
		return "HTML 4.01"
	case strings.Contains(doctype, "html public \"-//w3c//dtd xhtml 1.0"):
		return "XHTML 1.0"
	case strings.Contains(doctype, "html"):
		return "HTML5"
	default:
		return "Unknown HTML version"
	}
}
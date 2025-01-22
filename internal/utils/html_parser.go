package utils

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

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

// DetectHTMLVersion identifies the HTML version of the page
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

// CountHeadings counts the headings in the HTML
func CountHeadings(htmlContent string) map[string]int {
	headings := make(map[string]int)
	doc, _ := html.Parse(strings.NewReader(htmlContent))

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && strings.HasPrefix(n.Data, "h") && len(n.Data) == 2 {
			headings[n.Data]++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	return headings
}

// CountLinksConcurrently analyzes links concurrently
func CountLinksConcurrently(baseURL, htmlContent string) (int, int, int, error) {
	// internalChan := make(chan int)
	// externalChan := make(chan int)
	// inaccessibleChan := make(chan int)
	// errorChan := make(chan error)

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return 0, 0, 0, err
	}

	base, _ := url.Parse(baseURL)
	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	var wg sync.WaitGroup
	internal, external, inaccessible := 0, 0, 0

	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			parsedLink, err := url.Parse(link)
			if err != nil || parsedLink.Scheme == "" {
				return
			}

			if parsedLink.Host == base.Host {
				internal++
			} else {
				external++
			}

			resp, err := http.Head(link)
			if err != nil || resp.StatusCode >= 400 {
				inaccessible++
			}
		}(link)
	}

	wg.Wait()
	return internal, external, inaccessible, nil
}

// ContainsLoginForm detects login forms in the HTML
func ContainsLoginForm(htmlContent string) bool {
	return strings.Contains(htmlContent, "type=\"password\"")
}

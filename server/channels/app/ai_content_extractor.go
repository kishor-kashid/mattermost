// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"bytes"
	"net/url"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

// extractContent performs a lightweight extraction of title/description/text.
func (a *App) extractContent(rawURL string, body []byte, contentType string) LinkContent {
	title, desc, text := simpleHTMLExtract(body)

	u, _ := url.Parse(rawURL)
	domain := ""
	if u != nil {
		domain = u.Hostname()
	}

	contentTypeLabel := classifyContent(domain, contentType)
	readingTime := estimateReadingTime(text)

	return LinkContent{
		Title:       title,
		Description: desc,
		Text:        text,
		ContentType: contentTypeLabel,
		ReadingTime: readingTime,
		Domain:      domain,
		FaviconURL:  buildFaviconURL(u),
	}
}

func classifyContent(domain, contentType string) string {
	d := strings.ToLower(domain)
	switch {
	case strings.Contains(d, "github.com"):
		return "github"
	case strings.Contains(d, "stackoverflow.com"):
		return "stackoverflow"
	default:
		if strings.Contains(strings.ToLower(contentType), "html") {
			return "article"
		}
		return "unknown"
	}
}

func buildFaviconURL(u *url.URL) string {
	if u == nil {
		return ""
	}
	return u.Scheme + "://" + u.Host + "/favicon.ico"
}

// simpleHTMLExtract pulls <title>, meta description, and plain text from body.
func simpleHTMLExtract(body []byte) (string, string, string) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return "", "", ""
	}

	var title string
	var desc string
	var textBuilder strings.Builder

	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				title = strings.TrimSpace(n.FirstChild.Data)
			}
			if n.Data == "meta" {
				var name, content string
				for _, attr := range n.Attr {
					if strings.EqualFold(attr.Key, "name") && (strings.EqualFold(attr.Val, "description") || strings.EqualFold(attr.Val, "og:description")) {
						name = attr.Val
					}
					if strings.EqualFold(attr.Key, "content") {
						content = attr.Val
					}
				}
				if name != "" && content != "" && desc == "" {
					desc = strings.TrimSpace(content)
				}
			}
		}

		if n.Type == html.TextNode {
			textBuilder.WriteString(n.Data)
			textBuilder.WriteString(" ")
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}

	walk(doc)

	text := collapseWhitespace(textBuilder.String())
	if len(text) > 8000 {
		text = text[:8000]
	}

	return title, desc, text
}

func collapseWhitespace(s string) string {
	var b strings.Builder
	prevSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
			}
			prevSpace = true
		} else {
			b.WriteRune(r)
			prevSpace = false
		}
	}
	return strings.TrimSpace(b.String())
}

func estimateReadingTime(text string) int {
	words := strings.Fields(text)
	if len(words) == 0 {
		return 0
	}
	wpm := 200
	minutes := (len(words) + wpm - 1) / wpm
	if minutes == 0 {
		minutes = 1
	}
	return minutes
}

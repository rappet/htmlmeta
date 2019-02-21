package htmlmeta

import (
	"io"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// HTMLAnalyzer extracts metadata from an HTML document
type HTMLAnalyzer struct {
	BaseURL     *url.URL
	ConvertURLs bool
}

func extractAttr(attrs []html.Attribute, key string) string {
	for _, a := range attrs {
		if key == a.Key {
			return a.Val
		}
	}
	return ""
}

func stringContentOfNode(node *html.Node) string {
	var text string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		addText := ""
		if c.Type == html.TextNode {
			addText = c.Data
		} else if c.Type == html.ElementNode {
			addText = stringContentOfNode(c)
		}
		if text == "" {
			text = strings.TrimSpace(addText)
		} else {
			text = strings.TrimSpace(text) + " " + strings.TrimSpace(addText)
		}
	}
	return text
}

func dir(path string) string {
	if len(path) == 0 {
		return ""
	}

	i := len(path) - 1

	for i > 0 && '/' != path[i] {
		i--
	}
	return path[0:i]
}

func (analyzer *HTMLAnalyzer) makeURLAbsolute(u *url.URL) error {
	if analyzer.ConvertURLs && analyzer.BaseURL == nil {
		return ErrorBaseURLNotSet
	} else if analyzer.ConvertURLs {
		if !u.IsAbs() {
			u.Scheme = analyzer.BaseURL.Scheme
			u.User = analyzer.BaseURL.User
			u.Host = analyzer.BaseURL.Host
			u.Path = dir(analyzer.BaseURL.Path) + "/" + u.Path
		}
		return nil
	} else {
		return nil
	}
}

// GetPageMeta reads an HTML file from a reader and generates a PageMeta struct
func (analyzer *HTMLAnalyzer) GetPageMeta(r io.Reader) (*PageMeta, error) {
	pageMeta := &PageMeta{
		Links:  make([]LinkMeta, 0, 10),
		Images: make([]ImageMeta, 0, 10),
	}

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			href := extractAttr(n.Attr, "href")
			text := stringContentOfNode(n)

			parsedURL, _ := url.Parse(href)
			analyzer.makeURLAbsolute(parsedURL)
			if parsedURL != nil {
				pageMeta.Links = append(pageMeta.Links, LinkMeta{
					URL:  *parsedURL,
					Text: text,
				})
			}
		} else if n.Type == html.ElementNode && n.Data == "img" {
			src := extractAttr(n.Attr, "src")
			alt := extractAttr(n.Attr, "alt")
			width, _ := strconv.Atoi(extractAttr(n.Attr, "width"))
			height, _ := strconv.Atoi(extractAttr(n.Attr, "height"))

			parsedSrc, _ := url.Parse(src)
			if parsedSrc != nil {
				pageMeta.Images = append(pageMeta.Images, ImageMeta{
					Source:        *parsedSrc,
					AlternateText: alt,
					Width:         width,
					Height:        height,
				})
			}
		} else if n.Type == html.ElementNode && n.Data == "title" {
			pageMeta.Title = stringContentOfNode(n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return pageMeta, nil
}

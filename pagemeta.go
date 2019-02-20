package htmlmeta

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// PageMeta contains all metadata extracted from an html page
type PageMeta struct {
	Title  string
	Links  []LinkMeta
	Images []ImageMeta
}

// CreatePageMeta reads an HTML file from a reader and generates a PageMeta struct
func CreatePageMeta(r io.Reader) (*PageMeta, error) {
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
			href := ""
			text := ""
			for _, a := range n.Attr {
				if a.Key == "href" {
					href = a.Val
				}
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					if text == "" {
						text = strings.TrimSpace(c.Data)
					} else {
						text = text + " " + strings.TrimSpace(c.Data)
					}
				}
			}

			parsedURL, _ := url.Parse(href)
			if parsedURL != nil {
				pageMeta.Links = append(pageMeta.Links, LinkMeta{
					URL:  *parsedURL,
					Text: text,
				})
			}
		} else if n.Type == html.ElementNode && n.Data == "img" {
			src := ""
			alt := ""
			for _, a := range n.Attr {
				if a.Key == "src" {
					src = a.Val
				} else if a.Key == "alt" {
					alt = a.Val
				}
			}
			parsedSrc, _ := url.Parse(src)
			if parsedSrc != nil {
				pageMeta.Images = append(pageMeta.Images, ImageMeta{
					Source:        *parsedSrc,
					AlternateText: alt,
				})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return pageMeta, nil
}

// LinkMeta contains extracted metadata from an a tag
type LinkMeta struct {
	URL  url.URL
	Text string
}

func (meta LinkMeta) String() string {
	return fmt.Sprintf("%s,'%s'", meta.URL.String(), meta.Text)
}

// ImageMeta contains extracted metadata from an img tag
type ImageMeta struct {
	Source        url.URL
	AlternateText string
	Width         int
	Height        int
}

func (meta ImageMeta) String() string {
	return fmt.Sprintf("%s,'%s',(%dx%d)", meta.Source.String(), meta.AlternateText, meta.Width, meta.Height)
}

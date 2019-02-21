package htmlmeta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// PageMeta contains all metadata extracted from an html page
type PageMeta struct {
	Title  string      `json:"title,omitempty"`
	Links  []LinkMeta  `json:"links,omitempty"`
	Images []ImageMeta `json:"images,omitempty"`
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
			href := extractAttr(n.Attr, "href")
			text := stringContentOfNode(n)

			parsedURL, _ := url.Parse(href)
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

// LinkMeta contains extracted metadata from an a tag
type LinkMeta struct {
	URL  url.URL
	Text string
}

func (meta LinkMeta) String() string {
	return fmt.Sprintf("%s,'%s'", meta.URL.String(), meta.Text)
}

type basicLinkMeta struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

// MarshalJSON marshals LinkMeta as JSON. URL is saved as string.
func (meta LinkMeta) MarshalJSON() ([]byte, error) {
	linkMeta := basicLinkMeta{
		meta.URL.String(),
		meta.Text,
	}
	return json.Marshal(linkMeta)
}

// UnmarshalJSON unmarshals LinkMeta from JSON. Source is converted from string.
func (meta *LinkMeta) UnmarshalJSON(j []byte) error {
	var linkMeta basicLinkMeta

	err := json.Unmarshal(j, &linkMeta)
	if err != nil {
		return nil
	}

	url, err := url.Parse(linkMeta.URL)
	if err != nil {
		return nil
	}

	meta.URL = *url
	meta.Text = linkMeta.Text

	return nil
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

type basicImageMeta struct {
	Source        string `json:"src"`
	AlternateText string `json:"alt,omitempty"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
}

// MarshalJSON marshals ImageMeta as JSON. Source is saved as string.
func (meta ImageMeta) MarshalJSON() ([]byte, error) {
	imageMeta := basicImageMeta{
		Source:        meta.Source.String(),
		AlternateText: meta.AlternateText,
		Width:         meta.Width,
		Height:        meta.Height,
	}
	return json.Marshal(imageMeta)
}

// UnmarshalJSON unmarshals ImageMeta from JSON. Source is converted from string.
func (meta *ImageMeta) UnmarshalJSON(j []byte) error {
	var imageMeta basicImageMeta

	err := json.Unmarshal(j, &imageMeta)
	if err != nil {
		return nil
	}

	source, err := url.Parse(imageMeta.Source)
	if err != nil {
		return nil
	}

	meta.Source = *source
	meta.AlternateText = imageMeta.AlternateText
	meta.Width = imageMeta.Width
	meta.Height = imageMeta.Height

	return nil
}

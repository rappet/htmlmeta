package htmlmeta

import "net/url"

// PageMeta contains all metadata extracted from an html page
type PageMeta struct {
	Title  string
	Links  []url.URL
	Images []url.URL
}

// LinkMeta contains extracted metadata from an a tag
type LinkMeta struct {
	URL  url.URL
	Text string
}

// ImageMeta contains extracted metadata from an img tag
type ImageMeta struct {
	Source        url.URL
	AlternateText string
	Width         int
	Height        int
}

package htmlmeta

// PageMeta contains all metadata extracted from an html page
type PageMeta struct {
	Title  string      `json:"title,omitempty"`
	Links  []LinkMeta  `json:"links,omitempty"`
	Images []ImageMeta `json:"images,omitempty"`
}

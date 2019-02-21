package htmlmeta

import (
	"encoding/json"
	"fmt"
	"net/url"
)

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

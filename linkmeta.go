package htmlmeta

import (
	"encoding/json"
	"fmt"
	"net/url"
)

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

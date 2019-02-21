package htmlmeta

import (
	"encoding/json"
	"reflect"
	"testing"
)

var imagetests = []struct {
	imageMeta ImageMeta
	asString  string
	asJSON    string
}{
	{
		ImageMeta{parseURL("http://example.org/foo.png"), "example", 123, 456},
		"http://example.org/foo.png,'example',(123x456)",
		`{"src":"http://example.org/foo.png","alt":"example","width":123,"height":456}`,
	},
}

func TestImageMetaString(t *testing.T) {
	for _, tt := range imagetests {
		asString := tt.imageMeta.String()
		if asString != tt.asString {
			t.Errorf("got %s, want %s", asString, tt.asString)
		}
	}
}

func TestImageMetaJSONMarshal(t *testing.T) {
	for _, tt := range imagetests {
		asJSON, err := json.Marshal(tt.imageMeta)
		if err != nil {
			t.Errorf("could not marshal %v: %s", tt.imageMeta, err.Error())
		}

		if string(asJSON) != tt.asJSON {
			t.Errorf("got %s, want %s", string(asJSON), tt.asJSON)
		}
	}
}

func TestImageMetaJSONUnmarshal(t *testing.T) {
	for _, tt := range imagetests {
		var fromJSON ImageMeta
		err := json.Unmarshal([]byte(tt.asJSON), &fromJSON)
		if err != nil {
			t.Errorf("could not unmarshal %s: %s", tt.asJSON, err.Error())
		}

		if !reflect.DeepEqual(fromJSON, tt.imageMeta) {
			t.Errorf("got %v, want %v", fromJSON, tt.imageMeta)
		}
	}
}

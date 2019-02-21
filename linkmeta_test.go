package htmlmeta

import (
	"encoding/json"
	"reflect"
	"testing"
)

var linktests = []struct {
	linkMeta *LinkMeta
	asString string
	asJSON   string
}{
	{
		&LinkMeta{parseURL("http://example.org/"), "example"},
		"http://example.org/,'example'",
		`{"url":"http://example.org/","text":"example"}`,
	},
}

func TestLinkMetaString(t *testing.T) {
	for _, tt := range linktests {
		if tt.linkMeta != nil {
			asString := tt.linkMeta.String()
			if asString != tt.asString {
				t.Errorf("got %s, want %s", asString, tt.asString)
			}
		}
	}
}

func TestLinkMetaJSONMarshal(t *testing.T) {
	for _, tt := range linktests {
		if tt.linkMeta != nil {
			asJSON, err := json.Marshal(tt.linkMeta)
			if err != nil {
				t.Errorf("could not marshal %v: %s", tt.linkMeta, err.Error())
			}

			if string(asJSON) != tt.asJSON {
				t.Errorf("got %s, want %s", string(asJSON), tt.asJSON)
			}
		}
	}
}

func TestLinkMetaJSONUmarshal(t *testing.T) {
	for _, tt := range linktests {
		var fromJSON LinkMeta
		err := json.Unmarshal([]byte(tt.asJSON), &fromJSON)
		if tt.linkMeta == nil {
			if err == nil {
				t.Error("expected error")
			}
		} else {
			if err != nil {
				t.Errorf("could not unmarshal %s: %s", tt.asJSON, err.Error())
			}

			if !reflect.DeepEqual(fromJSON, *tt.linkMeta) {
				t.Errorf("got %v, want %v", fromJSON, tt.linkMeta)
			}
		}
	}
}

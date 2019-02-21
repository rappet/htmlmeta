package htmlmeta

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func parse_url(s string) url.URL {
	u, _ := url.Parse(s)
	return *u
}

var createMetaTests = []struct {
	in  string
	out *PageMeta
}{
	{
		`
			<!DOCTYPE html>
			<html>
			<body>
				<a href="http://example.org/foo">foo</a>
				<img src="http://example.org/a.png" alt="bar"/>
			</body>
			</html>
		`,
		&PageMeta{
			Links: []LinkMeta{
				LinkMeta{parse_url("http://example.org/foo"), "foo"},
			},
			Images: []ImageMeta{
				ImageMeta{parse_url("http://example.org/a.png"), "bar", 0, 0},
			},
		},
	},
	{
		`
			<body>
				<img src="http://example.org/size.jpg" alt="size" width="123" height="456"/>
			</body>
		`,
		&PageMeta{
			Links: []LinkMeta{},
			Images: []ImageMeta{
				ImageMeta{parse_url("http://example.org/size.jpg"), "size", 123, 456},
			},
		},
	},
	{
		`
			<body>
				<a href="http://example.org/baz">
					<img src="http://example.org/b.png" alt="foobar"/>
				</a>
			</body>
		`,
		&PageMeta{
			Links: []LinkMeta{
				LinkMeta{parse_url("http://example.org/baz"), ""},
			},
			Images: []ImageMeta{
				ImageMeta{parse_url("http://example.org/b.png"), "foobar", 0, 0},
			},
		},
	},
	{
		`<body><a href="http://example.com/">  Foo  <br/>  <p>Bar</p></a></body>`,
		&PageMeta{
			Links: []LinkMeta{
				LinkMeta{parse_url("http://example.com/"), "Foo Bar"},
			},
			Images: []ImageMeta{},
		},
	},
	{
		`<html><head><title>Foo</title></head></html>`,
		&PageMeta{
			Title:  "Foo",
			Links:  []LinkMeta{},
			Images: []ImageMeta{},
		},
	},
}

func TestCreatePageMeta(t *testing.T) {
	for _, tt := range createMetaTests {
		r := strings.NewReader(tt.in)
		meta, err := CreatePageMeta(r)
		if meta == nil && err != nil && tt.out == nil {
			continue
		}
		if !reflect.DeepEqual(meta, tt.out) {
			t.Logf("extracted: %v", meta)
			t.Logf("expected: %v", tt.out)
			t.Fatal("example meta and extracted meta are not equal")
		}
	}
}

func TestLinkMeta(t *testing.T) {
	linktests := []struct {
		linkMeta LinkMeta
		asString string
		asJSON   string
	}{
		{
			LinkMeta{parse_url("http://example.org/"), "example"},
			"http://example.org/,'example'",
			`{"url":"http://example.org/","text":"example"}`,
		},
	}
	for _, tt := range linktests {
		asString := tt.linkMeta.String()
		if asString != tt.asString {
			t.Errorf("got %s, want %s", asString, tt.asString)
		}

		asJSON, err := json.Marshal(tt.linkMeta)
		if err != nil {
			t.Errorf("could not marshal %v: %s", tt.linkMeta, err.Error())
		}

		if string(asJSON) != tt.asJSON {
			t.Errorf("got %s, want %s", string(asJSON), tt.asJSON)
		}

		var fromJSON LinkMeta
		err = json.Unmarshal([]byte(tt.asJSON), &fromJSON)
		if err != nil {
			t.Errorf("could not unmarshal %s: %s", tt.asJSON, err.Error())
		}

		if !reflect.DeepEqual(fromJSON, tt.linkMeta) {
			t.Errorf("got %v, want %v", fromJSON, tt.linkMeta)
		}
	}
}

func TestImageMeta(t *testing.T) {
	imagetests := []struct {
		imageMeta ImageMeta
		asString  string
		asJSON    string
	}{
		{
			ImageMeta{parse_url("http://example.org/foo.png"), "example", 123, 456},
			"http://example.org/foo.png,'example',(123x456)",
			`{"src":"http://example.org/foo.png","alt":"example","width":123,"height":456}`,
		},
	}
	for _, tt := range imagetests {
		asString := tt.imageMeta.String()
		if asString != tt.asString {
			t.Errorf("got %s, want %s", asString, tt.asString)
		}

		asJSON, err := json.Marshal(tt.imageMeta)
		if err != nil {
			t.Errorf("could not marshal %v: %s", tt.imageMeta, err.Error())
		}

		if string(asJSON) != tt.asJSON {
			t.Errorf("got %s, want %s", string(asJSON), tt.asJSON)
		}

		var fromJSON ImageMeta
		err = json.Unmarshal([]byte(tt.asJSON), &fromJSON)
		if err != nil {
			t.Errorf("could not unmarshal %s: %s", tt.asJSON, err.Error())
		}

		if !reflect.DeepEqual(fromJSON, tt.imageMeta) {
			t.Errorf("got %v, want %v", fromJSON, tt.imageMeta)
		}
	}
}

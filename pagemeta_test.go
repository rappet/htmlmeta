package htmlmeta

import (
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
}

func TestCreatePageMeta(t *testing.T) {
	for _, tt := range createMetaTests {
		r := strings.NewReader(tt.in)
		meta, err := CreatePageMeta(r)
		if err != nil {
			t.Fatal(err)
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
	}{
		{
			LinkMeta{parse_url("http://example.org/"), "example"},
			"http://example.org/,'example'",
		},
	}
	for _, tt := range linktests {
		asString := tt.linkMeta.String()
		if asString != tt.asString {
			t.Errorf("got %s, want %s", asString, tt.asString)
		}
	}
}

func TestImageMeta(t *testing.T) {
	imagetests := []struct {
		imageMeta ImageMeta
		asString  string
	}{
		{
			ImageMeta{parse_url("http://example.org/foo.png"), "example", 123, 456},
			"http://example.org/foo.png,'example',(123x456)",
		},
	}
	for _, tt := range imagetests {
		asString := tt.imageMeta.String()
		if asString != tt.asString {
			t.Errorf("got %s, want %s", asString, tt.asString)
		}
	}
}

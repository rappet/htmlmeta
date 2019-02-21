package htmlmeta

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func parseURL(s string) url.URL {
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
				LinkMeta{parseURL("http://example.org/foo"), "foo"},
			},
			Images: []ImageMeta{
				ImageMeta{parseURL("http://example.org/a.png"), "bar", 0, 0},
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
				ImageMeta{parseURL("http://example.org/size.jpg"), "size", 123, 456},
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
				LinkMeta{parseURL("http://example.org/baz"), ""},
			},
			Images: []ImageMeta{
				ImageMeta{parseURL("http://example.org/b.png"), "foobar", 0, 0},
			},
		},
	},
	{
		`<body><a href="http://example.com/">  Foo  <br/>  <p>Bar</p></a></body>`,
		&PageMeta{
			Links: []LinkMeta{
				LinkMeta{parseURL("http://example.com/"), "Foo Bar"},
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

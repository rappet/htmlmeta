package htmlmeta

import (
	"net/url"
	"reflect"
	"strings"
	"testing"
)

const exampleFile string = `
<!DOCTYPE html>
<html>
<body>
	<a href="http://example.org/foo">foo</a>
	<img src="http://example.org/a.png" alt="bar"/>
</body>
</html>
`

func parse_url(s string) url.URL {
	u, _ := url.Parse(s)
	return *u
}

var exampleMeta = &PageMeta{
	Title: "",
	Links: []LinkMeta{LinkMeta{
		URL:  parse_url("http://example.org/foo"),
		Text: "foo",
	}},
	Images: []ImageMeta{ImageMeta{
		Source:        parse_url("http://example.org/a.png"),
		AlternateText: "bar",
	}},
}

func TestCreatePageMeta(t *testing.T) {
	r := strings.NewReader(exampleFile)
	meta, err := CreatePageMeta(r)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(meta, exampleMeta) {
		t.Log("extracted:", meta)
		t.Log("expected:", exampleMeta)
		t.Fatal("example meta and extracted meta are not equal")
	}
}

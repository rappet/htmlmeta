package htmlmeta

import (
	"reflect"
	"testing"
)

func TestMakeURLAbsolute(t *testing.T) {
	absoluteTests := []struct {
		base string
		in   string
		out  string
	}{
		{
			"http://example.org",
			"foo.html",
			"http://example.org/foo.html",
		},
		{
			"https://example.org:8080/foo/index.html",
			"foo.html?1234#123",
			"https://example.org:8080/foo/foo.html?1234#123",
		},
	}

	for _, tt := range absoluteTests {
		base := parseURL(tt.base)
		analyzer := &HTMLAnalyzer{
			BaseURL:     &base,
			ConvertURLs: true,
		}

		in := parseURL(tt.in)
		out := parseURL(tt.out)
		err := analyzer.makeURLAbsolute(&in)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(in, out) {
			t.Errorf("got %v, want %v", in, tt.out)
		}
	}

	dummy := parseURL("foo.html")
	analyzer := &HTMLAnalyzer{
		ConvertURLs: true,
	}
	err := analyzer.makeURLAbsolute(&dummy)
	if err == nil {
		t.Error("expected error")
	}
}

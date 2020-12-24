package strings_test

import (
	"testing"

	"github.com/sinmetalcraft/ironhead/strings"
)

func TestReadNewLine(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		wantText string
		wantNext string
	}{
		{"シンプルに1改行", "hoge\nfuga", "hoge", "fuga"},
		{"末尾に改行", "hoge\n", "hoge", ""},
		{"改行がない", "hoge", "hoge", ""},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, next := strings.ReadNewLine(tt.text)
			if got != tt.wantText {
				t.Errorf("want text: %s but got %s", tt.wantText, got)
			}
			if next != tt.wantNext {
				t.Errorf("want continued: %s but got %s", tt.wantNext, next)
			}
		})
	}
}

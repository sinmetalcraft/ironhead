package erroreportings_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sinmetalcraft/ironhead/erroreportings"
)

func TestParse(t *testing.T) {
	mailBody := `
Google Cloud Platform
New error in xxxxx




context canceled

panic: context canceled goroutine 1 [running]: main.main()  
/layers/google.go.appengine_gopath/gopath/src/github.com/hoge/fuga/server/app.go:102  
+0x97c

main


Service

default

Version

v20201215b

View error details


If you no longer wish to receive messages like this one, you can  
unsubscribe.
`

	cases := []struct {
		name string
		text string
		want *erroreportings.Info
	}{
		{"Error Reporting Mail Parse",
			mailBody,
			&erroreportings.Info{
				ProjectID:      "xxxxx",
				Service:        "default",
				Version:        "v20201215b",
				ErrorDetailURL: "",
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := erroreportings.Parse(tt.text)
			if !cmp.Equal(got, tt.want) {
				t.Log(cmp.Diff(got, tt.want))
			}
		})
	}
}

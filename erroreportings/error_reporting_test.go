package erroreportings_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sinmetalcraft/ironhead/erroreportings"
)

func TestParse(t *testing.T) {
	plainTextMailBody := `
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

	htmlMailBody := `
<div style="font:13px Roboto, Verdana, sans-serif;padding:10px;color:#424242"><table style="padding:0;width:600px;background-color:#fafafa;margin:auto;border:1px solid #eeeeee;border-collapse:collapse;font:13px Roboto, Verdana, sans-serif"><tr><td style="padding:32px;background-color:#ededed;border-bottom:1px solid #e7e7e7"><a href="https://console.cloud.google.com?utm_source=cloud-notification&amp;utm_medium=email&amp;utm_content=new-error"><img src="https://cloud.google.com/_static/images/new-gcp-logo.png" width="236" height="30" alt="Google Cloud Platform"></a><div style="color:#8d8d8d;font-size:15px;padding-top:14px">New error in xxxxx</div></td></tr><tr><td style="padding:32px"><div style="padding-bottom:32px;font-size:18px"><div>context canceled</div><div>panic: context canceled

goroutine 1 [running]:
main.main()
	/layers/google.go.appengine_gopath/gopath/src/github.com/hoge/hogeWeb/server/app.go:102 +0x97c</div><div style="font-size:12px">main</div></div><div style="padding-bottom:32px"><div style="font-weight:bold">Service</div><div style="margin-bottom:16px">default</div><div style="font-weight:bold">Version</div><div style="margin-bottom:16px">v20201215b</div></div><a href="https://console.cloud.google.com/errors/CI3Gk9eWh-Wr3AE?project=xxxxx&amp;time=P30D&amp;utm_source=cloud-notification&amp;utm_medium=email&amp;utm_content=new-error" style="text-decoration:none"><span style="display:inline-block;background-color:#3271ed;color:white;font-size:13px;text-align:center;padding:12px 20px;border-radius:3px;line-height:16px;text-transform:uppercase">View error details</span></a></td></tr></table><div style="width:600px;margin:10px auto;text-align:center;font-size:10px">If you no longer wish to receive messages like this one, you can <a href="https://console.cloud.google.com/user-preferences/communication?project=xxxxx&amp;utm_source=cloud-notification&amp;utm_medium=email&amp;utm_content=unsubscribe-new-error" style="color:#3775EA;text-decoration:none"><b>unsubscribe</b></a>.</div></div>
`

	cases := []struct {
		name      string
		plainText string
		html      string
		want      *erroreportings.Info
	}{
		{"Error Reporting Mail Parse",
			plainTextMailBody,
			htmlMailBody,
			&erroreportings.Info{
				ProjectID:      "xxxxx",
				Service:        "default",
				Version:        "v20201215b",
				ErrorDetailURL: "https://console.cloud.google.com/errors/CI3Gk9eWh-Wr3AE?project=xxxxx&amp;time=P30D&amp;utm_source=cloud-notification&amp;utm_medium=email&amp;utm_content=new-error",
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got := erroreportings.Parse(tt.plainText, tt.html)
			if !cmp.Equal(got, tt.want) {
				t.Log(cmp.Diff(got, tt.want))
			}
		})
	}
}

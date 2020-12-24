package erroreportings

import (
	"regexp"
	"strings"

	ironstrings "github.com/sinmetalcraft/ironhead/strings"
)

// Info is ErrorReportingから送られてきたメールからぱっと必要な情報を抜き出したもの
type Info struct {
	ProjectID      string
	Service        string
	Version        string // optional
	ErrorDetailURL string
}

var consoleURLReg *regexp.Regexp

func init() {
	consoleURLReg = regexp.MustCompile(`https://console.cloud.google.com/errors/.*"`)
}

// Parse is ErrorReportingから送れてきたメールからぱっと必要な情報を抜き出す
// 引数にはErrorReportingのメールのplainTextを渡す
func Parse(plainText string, html string) *Info {
	info := &Info{}
	pID, next := readProjectID(plainText)
	info.ProjectID = pID

	service, next := readService(next)
	info.Service = service

	version, _ := readVersion(next)
	info.Version = version

	consoleURL := readConsoleURL(html)
	info.ErrorDetailURL = consoleURL

	return info
}

func readProjectID(text string) (string, string) {
	v := text[len("Google Cloud Platform\nNew error in "):]
	line, next := ironstrings.ReadNewLine(v)
	return strings.TrimSpace(line), next
}

func readService(text string) (string, string) {
	return readLabelText("Service", text)
}

func readVersion(text string) (string, string) {
	return readLabelText("Version", text)
}

// readLabelText is 指定したラベルが出てきた次の空行じゃない文字列を返す
func readLabelText(label string, text string) (string, string) {
	var next = text
	for {
		var line string
		line, next = ironstrings.ReadNewLine(next)
		if strings.TrimSpace(line) == label {
			break
		}
		if len(next) < 1 {
			return "", ""
		}
	}

	// labelで指定した文字列が出てきた後、最初の空行ではない行を返す
	for {
		var line string
		line, next = ironstrings.ReadNewLine(next)
		if len(line) > 0 {
			return strings.TrimSpace(line), next
		}
		if len(next) < 1 {
			return "", ""
		}
	}
}

func readConsoleURL(text string) string {
	v := consoleURLReg.FindString(text)
	var buf strings.Builder
	for _, c := range v {
		if string(c) == "\"" {
			return buf.String()
		}
		buf.WriteRune(c)
	}
	return buf.String()
}

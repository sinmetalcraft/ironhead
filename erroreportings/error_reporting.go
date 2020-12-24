package erroreportings

import (
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

// Parse is ErrorReportingから送れてきたメールからぱっと必要な情報を抜き出す
// 引数にはErrorReportingのメールのplainTextを渡す
func Parse(text string) *Info {
	info := &Info{}
	pID, next := readProjectID(text)
	info.ProjectID = pID

	service, next := readService(next)
	info.Service = service

	version, _ := readVersion(next)
	info.Version = version

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

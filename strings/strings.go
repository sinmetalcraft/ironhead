package strings

import "strings"

// ReadNewLine is 改行まで読んでいく。
// results : 読み取った文字列, 残った文字列
// 空行はSkipしていく
func ReadNewLine(text string) (string, string) {
	var buf strings.Builder
	for i := 0; i < len(text); i++ {
		if string(text[i]) == "\n" {
			if i == len(text)-1 {
				break
			}
			return buf.String(), text[i+1:]
		}
		buf.WriteRune(rune(text[i]))
	}
	return buf.String(), ""
}

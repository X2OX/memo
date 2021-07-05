package util

import (
	"strings"
)

var escapeMarkdownV2 = []string{
	"_", "*", "[", "]", "(", ")", "~", "`", ">",
	"#", "+", "-", "=", "|", "{", "}", ".", "!",
}
var escapeMarkdownV1 = []string{"_", "*", "`", "["}
var escapeHTML = [][]string{
	{"&", "&amp;"},
	{"<", "&lt;"},
	{">", "&gt;"},
}

func EscapedMarkdownV2(s string) string {
	for _, v := range escapeMarkdownV2 {
		s = strings.ReplaceAll(s, v, `\`+v)
	}
	return s
}
func EscapedMarkdownV1(s string) string {
	for _, v := range escapeMarkdownV1 {
		s = strings.ReplaceAll(s, v, `\`+v)
	}
	return s
}
func EscapedHTML(s string) string {
	for _, v := range escapeHTML {
		s = strings.ReplaceAll(s, v[0], v[1])
	}
	return s
}

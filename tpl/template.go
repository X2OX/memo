package tpl

import (
	"bytes"
	"html/template"
	"path/filepath"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/pkg/util"
)

func Load() *template.Template {
	if util.Exists(filepath.Join(model.Conf.TemplatesFolder(), "tpl.html")) {
		tpl, err := template.ParseFiles()
		if err == nil {
			return tpl
		}
	}

	tpl, _ := template.New("tpl.html").Parse(defaultHTML)
	return tpl
}

func ToHTML(s string) template.HTML {
	var buf bytes.Buffer

	if err := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	).Convert([]byte(s), &buf); err != nil {
		return ""
	}

	return template.HTML(buf.String())
}

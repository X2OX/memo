package model

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/x2ox/memo/pkg/dandelion"
	"github.com/x2ox/memo/pkg/participle"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"gorm.io/gorm"

	"github.com/x2ox/memo/pkg/util"
)

type Note struct {
	ID        uint64         `gorm:"primaryKey" json:"id" `
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `json:"title"`   // 标题
	Content   string         `json:"content"` // 内容
}

func (n *Note) ParticipleTitle() string   { return participle.Parse(n.Title) }
func (n *Note) ParticipleContent() string { return participle.Parse(n.Content) }
func (n *Note) link(t Type) string {
	return fmt.Sprintf("%s/preview/%s", Conf.Domain, NewToken(t, n.ID))
}
func (n *Note) PreviewLink() string { return n.link(Preview) }
func (n *Note) ViewLink() string    { return n.link(View) }
func (n *Note) ShareLink() string   { return n.link(Share) }
func (n *Note) List() string {
	return fmt.Sprintf(`%d \| %s%s%s \| %s
`,
		n.ID,
		"`", n.CreatedAt.Format("2006-01-02 15:04"), "`",
		n.MarkdownLink(),
	)
}
func (n *Note) MarkdownLink() string {
	return fmt.Sprintf(`[%s](%s)`, util.EscapedMarkdownV2(n.Title), n.ViewLink())
}
func (n *Note) HTML() template.HTML {
	if n == nil {
		return ""
	}
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
	).Convert([]byte(
		fmt.Sprintf("# %s\n\n %s \n\n <hr />%s", n.Title, n.Content, n.CreatedAt.Format("2006-01-02 15:04")),
	), &buf); err != nil {
		return ""
	}

	return template.HTML(buf.String())
}

func (n *Note) Description() string {
	if len([]rune(n.Content)) <= 256 {
		return n.Content
	}
	return string([]rune(n.Content)[:256]) + "..."
}

func (n *Note) InlineQueryResultArticle() dandelion.InlineQueryResultArticle {
	return dandelion.InlineQueryResultArticle{
		Type:  "article",
		ID:    strconv.FormatUint(n.ID, 10),
		Title: n.Title,
		InputMessageContent: dandelion.InputTextMessageContent{
			Text:                  n.Description(),
			DisableWebPagePreview: true,
		},
		ReplyMarkup: &dandelion.InlineKeyboardMarkup{
			InlineKeyboard: [][]dandelion.InlineKeyboardButton{
				{dandelion.NewInlineKeyboardButtonURL("查看内容", n.ShareLink())},
			},
		},
		Description: n.Description(),
		ThumbURL:    defaultImage,
	}
}

// NewNote 标题规则
// 如果设置了标题，默认为第一行作为标题，超过 32 个字符的，认定为无标题
func NewNote(text string) *Note {
	n := &Note{
		Title:   "无标题文档",
		Content: text,
	}
	str := strings.SplitN(text, "\n", 2)
	if len(str) == 2 && len([]rune(str[0])) <= 32 {
		n.Title = str[0]
		n.Content = str[1]
	}
	return n
}

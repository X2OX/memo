package model

import (
	"fmt"

	"github.com/x2ox/memo/pkg/dandelion"
)

const (
	defaultImage = `https://cdn.jsdelivr.net/gh/x2ox/note@f5ea0450547408cf661331f18b6da04120bae500/.data/static/android-chrome-192x192.png`
)

var (
	Username = `aoangc`
	header   = `📝[*Memo @*](https://t.me/%s)\#%s
`
	pagination = "\nPage: %d/%d  Count: %d"
)

func Header(tag string) string {
	return fmt.Sprintf(header, Username, tag)
}

func Pagination(current, total, count int64) string {
	return fmt.Sprintf(pagination, current, total, count)
}

func SwitchMode(mode string) string {
	return fmt.Sprintf("模式已切换到: %s", mode)
}

func NoMoreContent() dandelion.InlineQueryResultArticle {
	return dandelion.InlineQueryResultArticle{
		Type:  "article",
		ID:    "NoMoreContent",
		Title: "没有更多的内容了",
		InputMessageContent: dandelion.InputTextMessageContent{
			Text:                  "没有更多的内容了",
			DisableWebPagePreview: true,
		},
		ThumbURL: defaultImage,
	}
}

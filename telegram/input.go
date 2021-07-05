package telegram

import (
	"bytes"
	"path"
	"strconv"
	"sync"

	"github.com/x2ox/memo/db"
	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/pkg/dandelion"
	"github.com/x2ox/memo/pkg/participle"
)

var inputAdapter = &Message{mode: model.ModeInput}

type Message struct {
	mux  sync.RWMutex
	mode model.Mode
}

func (i *Message) Mode() model.Mode {
	i.mux.RLock()
	mode := model.ModeInput
	if i.mode == model.ModeSearch {
		mode = model.ModeSearch
	}
	i.mux.RUnlock()
	return mode
}

func (i *Message) SwitchMode() string {
	i.mux.Lock()
	if i.mode == model.ModeInput {
		i.mode = model.ModeSearch
	} else {
		i.mode = model.ModeInput
	}
	i.mux.Unlock()
	return i.Mode().String()
}

func (i *Message) Adapter() dandelion.Adapters { return nil }
func (i *Message) IsMatch(c *dandelion.Context) bool {
	return c.Message.Message != nil && !c.Message.Message.IsCommand()
}
func (i *Message) Handle(c *dandelion.Context) bool {
	if i.Mode() == model.ModeSearch {
		searchMode(c)
	} else {
		inputMode(c)
	}
	return true
}

type EditedMessage struct{}

func (EditedMessage) Adapter() dandelion.Adapters       { return nil }
func (EditedMessage) IsMatch(c *dandelion.Context) bool { return c.Message.EditedMessage != nil }
func (EditedMessage) Handle(c *dandelion.Context) bool {
	if c.Message.EditedMessage.Text != "" { // 只处理有文本的编辑
		if err := db.Input.UpdateContent(c.Message.EditedMessage.MessageID, c.Message.EditedMessage.Text); err != nil {
		}
	}
	return true
}

func searchMode(c *dandelion.Context) {
	if c.Message.Message.Text == "" {
		c.SendText("ヽ(*。>Д<)o゜ 只能搜索文本哦")
		return
	}

	notes, count := db.Search.Search(participle.Parse(c.Message.Message.Text), 0, 15)
	countPage := count / 15
	if count%15 != 0 {
		countPage++
	}

	ikb := make([]dandelion.InlineKeyboardButton, 0, 1)
	if countPage > 1 {
		ikb = append(ikb, dandelion.InlineKeyboardButton{
			Text:         "下一页",
			CallbackData: NewCallbackData(CallbackTypeSearch, c.Message.Message.Text, strconv.Itoa(2)),
		})
	}

	var buf bytes.Buffer
	buf.WriteString(model.Header("Search `" + c.Message.Message.Text + "`\n\n"))
	for _, v := range notes {
		buf.WriteString(v.List())
	}
	buf.WriteString(model.Pagination(1, countPage, count))

	_, _ = c.Send(c.NewMessage(
		buf.String(),
		&dandelion.InlineKeyboardMarkup{
			InlineKeyboard: [][]dandelion.InlineKeyboardButton{ikb},
		},
	))
}

func inputMode(c *dandelion.Context) {
	input := &model.Input{
		MessageID: c.Message.Message.MessageID,
	}
	var buf bytes.Buffer

	if c.Message.Message.Text != "" {
		buf.WriteString(c.Message.Message.Text)
		buf.WriteByte('\n')
	}
	if c.Message.Message.Animation != nil {
		if filename := c.DownloadAndSave(
			c.Message.Message.Animation.FileID, "/data/memo/file/"); filename != "" {
			buf.WriteString(linkMarkdown(filename))
			buf.WriteByte('\n')
		}
	}
	if len(c.Message.Message.Photo) > 0 {
		if filename := c.DownloadAndSave(
			c.Message.Message.Photo[len(c.Message.Message.Photo)-1].FileID,
			"/data/memo/file/"); filename != "" {
			buf.WriteString(linkMarkdown(filename))
			buf.WriteByte('\n')
		}
	}
	if c.Message.Message.Document != nil {
		if filename := c.DownloadAndSave(
			c.Message.Message.Document.FileID, "/data/memo/file/"); filename != "" {
			buf.WriteString(linkMarkdown(filename))
			buf.WriteByte('\n')
		}
	}
	if c.Message.Message.Video != nil {
		if filename := c.DownloadAndSave(
			c.Message.Message.Video.FileID, "/data/memo/file/"); filename != "" {
			buf.WriteString(linkMarkdown(filename))
			buf.WriteByte('\n')
		}
	}
	if c.Message.Message.Audio != nil {
		if filename := c.DownloadAndSave(
			c.Message.Message.Audio.FileID, "/data/memo/file/"); filename != "" {
			buf.WriteString(linkMarkdown(filename))
			buf.WriteByte('\n')
		}
	}
	if c.Message.Message.Voice != nil {
		if filename := c.DownloadAndSave(
			c.Message.Message.Voice.FileID, "/data/memo/file/"); filename != "" {
			buf.WriteString(linkMarkdown(filename))
			buf.WriteByte('\n')
		}
	}
	if c.Message.Message.VideoNote != nil {
		if filename := c.DownloadAndSave(
			c.Message.Message.VideoNote.FileID, "/data/memo/file/"); filename != "" {
			buf.WriteString(linkMarkdown(filename))
			buf.WriteByte('\n')
		}
	}

	if input.Content = buf.String(); input.Content != "" { // 跳过
		if err := db.Input.Add(input); err != nil {
		}
	}
}

var imgExt = []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp"}

func linkMarkdown(filename string) string {
	ext := path.Ext(filename)
	for _, v := range imgExt {
		if v == ext {
			return `![](/file/` + filename + ")\n"
		}
	}
	return `[` + filename + `](/file/` + filename + ")\n"
}

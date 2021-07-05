package dandelion

import (
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Context struct {
	Engine  *Engine
	Message Update

	ID string
}

func (c *Context) reset() {

}

func (c *Context) SendText(s string) {
	_, err := c.Engine.Send(NewMessage(c.Message.Message.From.ID, s))
	if err != nil {
		return
	}
}
func (c *Context) ReplyText(s string) {
	_, err := c.Engine.Send(MessageConfig{
		BaseChat: BaseChat{
			ChatID:           c.Message.Message.Chat.ID,
			ReplyToMessageID: c.Message.Message.MessageID,
		},
		Text:      s,
		ParseMode: ModeMarkdownV2,
	})
	if err != nil {
		return
	}
}

func (c *Context) Send(ct Chattable) (Message, error) {
	return c.Engine.Send(ct)
}

func (c *Context) DownloadAndSave(fileID, fileDir string) string {
	var (
		resp    *http.Response
		err     error
		file    *os.File
		fileUrl = c.GetFileDirectURL(fileID)
	)
	if fileUrl == "" {
		return ""
	}

	if resp, err = http.Get(fileUrl); err != nil {
		return ""
	}
	defer resp.Body.Close()

	filename := fileID + path.Ext(fileUrl)
	if file, err = os.OpenFile(filepath.Join(fileDir, filename),
		os.O_RDWR|os.O_CREATE|os.O_TRUNC, 666); err != nil {
		return ""
	}
	_, err = io.Copy(file, resp.Body)
	return filename
}

func (c *Context) GetFileDirectURL(fileID string) string {
	u, _ := c.Engine.GetFileDirectURL(fileID)
	return u
}
func (c *Context) IsMessageToMe(message Message) bool {
	return strings.Contains(message.Text, "@"+c.Engine.Self.UserName)
}

func (c *Context) CommandIs(command string) bool {
	return c != nil && c.Message.Message != nil && c.Message.Message.Command() == command
}
func (c *Context) NewEditListMessage(text string, ikb *InlineKeyboardMarkup) EditMessageTextConfig {
	if c.Message.CallbackQuery == nil {
		return EditMessageTextConfig{}
	}
	return EditMessageTextConfig{
		BaseEdit: BaseEdit{
			ChatID:      c.Message.CallbackQuery.Message.Chat.ID,
			MessageID:   c.Message.CallbackQuery.Message.MessageID,
			ReplyMarkup: ikb,
		},
		Text:      text,
		ParseMode: ModeMarkdownV2,
	}
}

func (c *Context) NewMessage(text string, ikb *InlineKeyboardMarkup) MessageConfig {
	if c.Message.Message == nil {
		return MessageConfig{}
	}
	return MessageConfig{
		BaseChat: BaseChat{
			ChatID:      c.Message.Message.Chat.ID,
			ReplyMarkup: ikb,
		},
		Text:                  text,
		ParseMode:             ModeMarkdownV2,
		Entities:              nil,
		DisableWebPagePreview: false,
	}
}

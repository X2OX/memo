package telegram

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/x2ox/memo/db"
	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/pkg/dandelion"
)

type CallbackDataType uint8

const (
	CallbackNone CallbackDataType = iota
	CallbackTypeSearch
	CallbackTypeList
	CallbackTypeUpdateKey
	CallbackTypeSetCommand
)

func NewCallbackData(t CallbackDataType, param ...string) *string {
	if b, err := json.Marshal(&CallbackData{
		Type:  t,
		Param: param,
	}); err == nil && len(b) > 0 {
		s := string(b)
		return &(s)
	}
	return nil
}
func (c CallbackData) Is(t CallbackDataType) bool { return c.Type == t }

func ParseCallbackData(s string) CallbackData {
	var arg CallbackData
	_ = json.Unmarshal([]byte(s), &arg)
	return arg
}

type CallbackData struct {
	Type  CallbackDataType `json:"t"`           // 数据类型
	Param []string         `json:"p,omitempty"` // 参数
}

type (
	Callback           struct{}
	CallbackSearch     struct{}
	CallbackList       struct{}
	CallbackUpdateKey  struct{}
	CallbackSetCommand struct{}
)

func (Callback) Adapter() dandelion.Adapters {
	return []dandelion.Adapter{
		&CallbackList{}, &CallbackSearch{}, &CallbackUpdateKey{},
		&CallbackSetCommand{},
	}
}
func (Callback) IsMatch(c *dandelion.Context) bool {
	return c.Message.CallbackQuery != nil && !ParseCallbackData(c.Message.CallbackQuery.Data).Is(CallbackNone)
}
func (Callback) Handle(*dandelion.Context) bool { return false }

func (CallbackSearch) Adapter() dandelion.Adapters { return nil }
func (CallbackSearch) IsMatch(c *dandelion.Context) bool {
	return ParseCallbackData(c.Message.CallbackQuery.Data).Is(CallbackTypeSearch)
}
func (CallbackSearch) Handle(c *dandelion.Context) bool {
	param := ParseCallbackData(c.Message.CallbackQuery.Data).Param
	if len(param) != 2 { // content string, page int
		return true
	}
	content := param[0]
	page, err := strconv.Atoi(param[1])
	if content == "" || err != nil || page <= 0 {
		return true
	}

	arr, count := db.Search.Search(content, (page-1)*15, 15)
	countPage := count / 15
	if count%15 != 0 {
		countPage++
	}

	ikb := make([]dandelion.InlineKeyboardButton, 0, 2)
	if page > 1 {
		ikb = append(ikb, dandelion.InlineKeyboardButton{
			Text:         "上一页",
			CallbackData: NewCallbackData(CallbackTypeSearch, strconv.Itoa(page-1)),
		})
	}
	if countPage > int64(page) {
		ikb = append(ikb, dandelion.InlineKeyboardButton{
			Text:         "下一页",
			CallbackData: NewCallbackData(CallbackTypeSearch, strconv.Itoa(page+1)),
		})
	}

	var buf bytes.Buffer
	buf.WriteString(model.Header("Search"))
	buf.WriteString(" `" + content + "`\n\n")
	for _, v := range arr {
		buf.WriteString(v.List())
	}

	buf.WriteString(model.Pagination(int64(page), countPage, count))

	_, _ = c.Send(c.NewEditListMessage(
		buf.String(),
		&dandelion.InlineKeyboardMarkup{
			InlineKeyboard: [][]dandelion.InlineKeyboardButton{ikb},
		},
	))

	return true
}

func (CallbackList) Adapter() dandelion.Adapters { return nil }
func (CallbackList) IsMatch(c *dandelion.Context) bool {
	return ParseCallbackData(c.Message.CallbackQuery.Data).Type == CallbackTypeList
}
func (CallbackList) Handle(c *dandelion.Context) bool {
	param := ParseCallbackData(c.Message.CallbackQuery.Data).Param
	if len(param) != 1 {
		return true
	}
	page, _ := strconv.Atoi(param[0])
	if page <= 0 {
		return true
	}

	arr, count := db.Note.Query((page-1)*15, 15)
	countPage := count / 15
	if count%15 != 0 {
		countPage++
	}

	ikb := make([]dandelion.InlineKeyboardButton, 0, 2)
	if page > 1 {
		ikb = append(ikb, dandelion.InlineKeyboardButton{
			Text:         "上一页",
			CallbackData: NewCallbackData(CallbackTypeList, strconv.Itoa(page-1)),
		})
	}
	if countPage > int64(page) {
		ikb = append(ikb, dandelion.InlineKeyboardButton{
			Text:         "下一页",
			CallbackData: NewCallbackData(CallbackTypeList, strconv.Itoa(page+1)),
		})
	}

	var buf bytes.Buffer
	buf.WriteString(model.Header("List"))
	buf.WriteString("\n\n")
	for _, v := range arr {
		buf.WriteString(v.List())
	}
	buf.WriteString(model.Pagination(int64(page), countPage, count))

	_, _ = c.Send(c.NewEditListMessage(
		buf.String(),
		&dandelion.InlineKeyboardMarkup{
			InlineKeyboard: [][]dandelion.InlineKeyboardButton{ikb},
		},
	))
	return true
}

func (CallbackUpdateKey) Adapter() dandelion.Adapters { return nil }
func (CallbackUpdateKey) IsMatch(c *dandelion.Context) bool {
	return ParseCallbackData(c.Message.CallbackQuery.Data).Type == CallbackTypeUpdateKey
}
func (CallbackUpdateKey) Handle(c *dandelion.Context) bool {
	model.UpdateKey()
	_, _ = c.Send(dandelion.CallbackConfig{
		CallbackQueryID: c.Message.CallbackQuery.ID,
		Text:            "重置密钥完成！",
	})
	return true
}

func (CallbackSetCommand) Adapter() dandelion.Adapters { return nil }
func (CallbackSetCommand) IsMatch(c *dandelion.Context) bool {
	return ParseCallbackData(c.Message.CallbackQuery.Data).Type == CallbackTypeSetCommand
}
func (CallbackSetCommand) Handle(c *dandelion.Context) bool {
	_, _ = c.Send(dandelion.SetMyCommandsConfig{}.Set([]dandelion.BotCommand{
		{Command: "start", Description: "「扬帆，起航！」"},
		{Command: "list", Description: "「随 手 笺」"},
		{Command: "mode", Description: "「响应模式」"},
		{Command: "preview", Description: "「预览草稿」"},
		{Command: "submit", Description: "「提交内容」"},
		{Command: "clear", Description: "「清空草稿」"},
	}))
	_, _ = c.Send(dandelion.CallbackConfig{
		CallbackQueryID: c.Message.CallbackQuery.ID,
		Text:            "同步完成，请等待生效",
	})
	return true
}

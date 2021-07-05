package telegram

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/x2ox/memo/db"
	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/pkg/dandelion"
	"go.uber.org/zap"
	"go.x2ox.com/blackdatura"
)

var (
	log    *zap.Logger
	engine *dandelion.Engine
)

func Init() {
	log = blackdatura.With("telegram")

	var err error
	if engine, err = dandelion.New(model.Conf.TelegramToken); err != nil {
		log.Fatal(
			"client telegram server error",
			zap.String("token", model.Conf.TelegramToken),
			zap.Error(err))
	}

	engine.SetAdapter(&Auth{})
	model.Username = engine.Username()

	if !model.Conf.IsWebhook() {
		engine.Run()
		return
	}
	if wi, _ := engine.GetWebhookInfo(); wi.URL == model.Conf.Webhook() {
		return
	}

	var u *url.URL
	if u, err = url.Parse(model.Conf.Webhook()); err != nil {
		log.Fatal("telegram webhook addr parse err", zap.Error(err))
	}
	if _, err = engine.SendSet(dandelion.WebhookConfig{URL: u}); err != nil {
		log.Fatal("set webhook error", zap.Error(err))
	}
}

func Close() {
	engine.Stop()
}

func Webhook() func(c *gin.Context) {
	return func(c *gin.Context) {
		engine.ListenForWebhook(c.Writer, c.Request)
	}
}

type (
	Command        struct{}
	CommandSubmit  struct{}
	CommandClear   struct{}
	CommandPreview struct{}
	CommandList    struct{}
	CommandMode    struct{}
	CommandStart   struct{}
	CommandDelete  struct{}
)

func (Command) Adapter() dandelion.Adapters {
	return []dandelion.Adapter{
		&CommandList{}, &CommandClear{}, &CommandSubmit{}, &CommandMode{},
		&CommandPreview{}, &CommandStart{}, &CommandDelete{},
	}
}
func (Command) IsMatch(c *dandelion.Context) bool {
	return c.Message.Message != nil && c.Message.Message.IsCommand()
}
func (Command) Handle(*dandelion.Context) bool { return false }

func (CommandSubmit) Adapter() dandelion.Adapters       { return nil }
func (CommandSubmit) IsMatch(c *dandelion.Context) bool { return c.CommandIs("submit") }
func (CommandSubmit) Handle(c *dandelion.Context) bool {
	if !db.Input.Check() {
		c.ReplyText(`ヽ\(\*。\>Д<\)o゜ 草稿箱内还是空的呢`)
		return true
	}
	note := db.Input.Submit()
	if note == nil {
		c.ReplyText(`\(；￣Д￣）似乎发生了点儿什么`)
		return true
	}

	c.ReplyText(`ฅ՞•ﻌ•՞ฅ 提交完成`)

	return true
}

func (CommandClear) Adapter() dandelion.Adapters       { return nil }
func (CommandClear) IsMatch(c *dandelion.Context) bool { return c.CommandIs("clear") }
func (CommandClear) Handle(c *dandelion.Context) bool {
	if err := db.Input.Clear(); err != nil {
	}
	c.ReplyText(`ヽ\(\*。\>Д<\)o゜ 草稿箱被清空了`)
	return true
}

func (CommandPreview) Adapter() dandelion.Adapters       { return nil }
func (CommandPreview) IsMatch(c *dandelion.Context) bool { return c.CommandIs("preview") }
func (CommandPreview) Handle(c *dandelion.Context) bool {
	u := fmt.Sprintf("%s/preview/%s", model.Conf.Domain, model.NewToken(model.Preview, 0))

	c.Send(c.NewMessage(
		fmt.Sprintf("%s\n 草稿箱内一共有: %d 条输入", model.Header("Preview"), db.Input.Count()),
		&dandelion.InlineKeyboardMarkup{
			InlineKeyboard: [][]dandelion.InlineKeyboardButton{{dandelion.InlineKeyboardButton{
				Text: "点击预览",
				URL:  &u,
			}}},
		},
	))
	return true
}

func (CommandList) Adapter() dandelion.Adapters       { return nil }
func (CommandList) IsMatch(c *dandelion.Context) bool { return c.CommandIs("list") }
func (CommandList) Handle(c *dandelion.Context) bool {
	arr, count := db.Note.Query(0, 15)
	countPage := count / 15
	if count%15 != 0 {
		countPage++
	}

	ikb := make([]dandelion.InlineKeyboardButton, 0, 1)
	if countPage > 1 {
		ikb = append(ikb, dandelion.InlineKeyboardButton{
			Text:         "下一页",
			CallbackData: NewCallbackData(CallbackTypeList, strconv.Itoa(2)),
		})
	}

	var buf bytes.Buffer
	buf.WriteString(model.Header("List"))
	buf.WriteString("\n\n")
	for _, v := range arr {
		buf.WriteString(v.List())
	}
	buf.WriteString(model.Pagination(1, countPage, count))

	_, _ = c.Send(c.NewMessage(
		buf.String(),
		&dandelion.InlineKeyboardMarkup{
			InlineKeyboard: [][]dandelion.InlineKeyboardButton{ikb},
		},
	))

	return true
}

func (CommandMode) Adapter() dandelion.Adapters       { return nil }
func (CommandMode) IsMatch(c *dandelion.Context) bool { return c.CommandIs("mode") }
func (CommandMode) Handle(c *dandelion.Context) bool {
	c.ReplyText(model.SwitchMode(inputAdapter.SwitchMode()))
	return true
}

func (CommandStart) Adapter() dandelion.Adapters       { return nil }
func (CommandStart) IsMatch(c *dandelion.Context) bool { return c.CommandIs("start") }
func (CommandStart) Handle(c *dandelion.Context) bool {
	_, _ = c.Send(dandelion.MessageConfig{
		BaseChat: dandelion.BaseChat{
			ChatID: c.Message.Message.Chat.ID,
			ReplyMarkup: &dandelion.InlineKeyboardMarkup{
				InlineKeyboard: [][]dandelion.InlineKeyboardButton{
					{
						dandelion.InlineKeyboardButton{
							Text:         "重置密钥",
							CallbackData: NewCallbackData(CallbackTypeUpdateKey),
						},
						dandelion.InlineKeyboardButton{
							Text:         "同步命令",
							CallbackData: NewCallbackData(CallbackTypeSetCommand),
						},
					},
				},
			},
		},
		Text: fmt.Sprintf("%s%s", model.Header("Start"), fmt.Sprintf(`
*船长*: 扬帆，起航！
*水手*: 什么方向？
*船长*: 海盗需要什么方向！起航！
*水手*: ……
*船长*: 拿笔记一下，免得回不来

`+"截至目前，一共有 `%d` 篇，最近一周新增 `%d` 篇，草稿箱内有 `%d` 条笔记",
			db.Note.Count(), db.Note.WeekCount(), db.Input.Count())),
		ParseMode: dandelion.ModeMarkdownV2,
	})
	return true
}

func (CommandDelete) Adapter() dandelion.Adapters       { return nil }
func (CommandDelete) IsMatch(c *dandelion.Context) bool { return c.CommandIs("delete") }
func (CommandDelete) Handle(c *dandelion.Context) bool {
	id, _ := strconv.ParseUint(c.Message.Message.CommandArguments(), 10, 64)
	if db.Note.Delete(id) != nil {
		c.ReplyText(`\(；￣Д￣）似乎那里不大对`)
		return true
	}
	c.ReplyText(`\(；￣Д￣）删除完成`)
	return true
}

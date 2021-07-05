package telegram

import (
	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/pkg/dandelion"
)

type Auth struct{}

func (Auth) Adapter() dandelion.Adapters {
	return []dandelion.Adapter{&Command{}, inputAdapter, &Inline{}, &Callback{}}
}
func (Auth) IsMatch(c *dandelion.Context) bool { return true }
func (Auth) Handle(c *dandelion.Context) bool {
	if c == nil ||
		(c.Message.Message != nil && c.Message.Message.From.ID != model.Conf.TelegramID) ||
		(c.Message.EditedMessage != nil && c.Message.EditedMessage.From.ID != model.Conf.TelegramID) ||
		(c.Message.InlineQuery != nil && c.Message.InlineQuery.From.ID != model.Conf.TelegramID) ||
		(c.Message.ChosenInlineResult != nil && c.Message.ChosenInlineResult.From.ID != model.Conf.TelegramID) ||
		(c.Message.CallbackQuery != nil && c.Message.CallbackQuery.From.ID != model.Conf.TelegramID) {
		return true
	}
	return false
}

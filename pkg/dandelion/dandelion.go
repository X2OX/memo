package dandelion

import (
	"context"
	"sync"
)

// Engine allows you to interact with the Telegram Bot API.
type Engine struct {
	Token string `json:"token"`
	Debug bool   `json:"debug"`

	Self   User       `json:"-"`
	Client HTTPClient `json:"-"`

	apiEndpoint  string
	UpdateConfig UpdateConfig
	ctx          context.Context
	cancelFunc   context.CancelFunc

	adapter Adapters
	pool    sync.Pool
}

func (bot *Engine) SetAdapter(adapters ...Adapter) {
	bot.adapter = adapters
}

func (bot *Engine) serve(msg Update) {
	go func() {
		c := bot.pool.Get().(*Context)
		defer func() {
			if err := recover(); err != nil {
			}
			bot.pool.Put(c)
		}()

		c.reset()
		c.Message = msg
		bot.adapter.Match(c)
	}()
}

func (bot *Engine) Username() string {
	return bot.Self.UserName
}

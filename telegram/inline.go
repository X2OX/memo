package telegram

import (
	"strconv"

	"github.com/x2ox/memo/db"
	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/pkg/dandelion"
	"github.com/x2ox/memo/pkg/participle"
	"go.uber.org/zap"
)

type Inline struct{}

func (Inline) Adapter() dandelion.Adapters       { return nil }
func (Inline) IsMatch(c *dandelion.Context) bool { return c.Message.InlineQuery != nil }
func (Inline) Handle(c *dandelion.Context) bool {
	var (
		notes     []*model.Note
		count     int64
		offset, _ = strconv.Atoi(c.Message.InlineQuery.Offset)
	)

	if c.Message.InlineQuery.Query != "" {
		notes, count = db.Search.Search(participle.Parse(c.Message.InlineQuery.Query), offset, 15)
	}

	arr := make([]interface{}, 0, len(notes))
	for _, v := range notes {
		arr = append(arr, v.InlineQueryResultArticle())
	}

	if len(arr) == 0 {
		arr = append(arr, model.NoMoreContent())
	}

	result := dandelion.InlineConfig{
		InlineQueryID: c.Message.InlineQuery.ID,
		Results:       arr,
		IsPersonal:    true,
	}

	if count != 0 && len(notes) >= 15 && int64(offset)+15 > count {
		result.NextOffset = strconv.Itoa(offset + 15)
	}

	if _, err := c.Send(result); err != nil {
		log.Warn("inline message send error", zap.Error(err))
	}
	return true
}

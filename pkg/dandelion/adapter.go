package dandelion

type Handler func(*Context) error

// Adapter 消息适配器，实现这个接口的适配器才行
type Adapter interface {
	Adapter() Adapters     // 子适配器，先运行子适配器
	IsMatch(*Context) bool // 是否匹配
	Handle(*Context) bool  // 是否结束
}

type Adapters []Adapter

func (a Adapters) Match(c *Context) bool {
	if len(a) == 0 {
		return true
	}

	for _, v := range a {
		if v.IsMatch(c) {
			if ok := v.Handle(c); ok {
				return false
			}
			return v.Adapter().Match(c)
		}
	}
	return false
}

package model

type Mode string

const (
	ModeInput  Mode = "输入模式"
	ModeSearch Mode = "搜索模式"
)

func (m Mode) String() string { return string(m) }

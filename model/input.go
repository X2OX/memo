package model

import (
	"time"

	"gorm.io/gorm"
)

type Input struct {
	ID        uint64         `gorm:"primaryKey" json:"id" `
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	MessageID int
	Content   string `json:"content"` // 内容
}

package db

import (
	"github.com/x2ox/memo/model"
	"gorm.io/gorm"
)

// FullTextSearch interface
type FullTextSearch interface {
	New(db *gorm.DB) FullTextSearch
	Init() error
	Index() error
	ReIndex() error
	Clean() error

	Search(keywords string, offset, limit int) ([]*model.Note, int64)
	Create(titleKeywords, contentKeywords string, id uint64) error
	Delete(id uint64) error
}

var Search FullTextSearch = &SQLite{}

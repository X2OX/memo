package db

import (
	"github.com/x2ox/memo/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initSQLite() gorm.Dialector {
	return sqlite.Open(model.Conf.DSN())
}

type SQLite struct{ db *gorm.DB }

func (s SQLite) Init() error {
	var count int64
	if err := s.db.Table("sqlite_master").Where("type = 'table'").
		Where("name = ?", "note_row").Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return s.db.Exec(`CREATE VIRTUAL TABLE note_row USING fts5(id UNINDEXED, title, content)`).Error
	}
	return nil
}

func (s SQLite) Index() error   { return nil }
func (s SQLite) ReIndex() error { return nil }
func (s SQLite) Clean() error   { return s.db.Exec("DROP TABLE note_row").Error }

func (s SQLite) Search(keywords string, offset, limit int) (arr []*model.Note, count int64) {
	s.db.Model(&model.Note{}).Where("id IN (?)", db.
		Select("id").Table("note_row").
		Where("note_row MATCH ?", keywords).
		Order("rank").
		Offset(offset).
		Limit(limit),
	).Find(&arr)
	s.db.Model(&model.Note{}).Where("id IN (?)", db.
		Select("id").Table("note_row").
		Where("note_row MATCH ?", keywords)).Count(&count)
	return
}
func (s SQLite) Create(titleKeywords, contentKeywords string, id uint64) error {
	return s.db.Exec(`INSERT INTO "note_row"("id", "title", "content") VALUES (?, ?, ?)`,
		id, titleKeywords, contentKeywords).Error
}
func (s SQLite) Delete(id uint64) error {
	return s.db.Exec("DELETE FROM `note_row` WHERE `note_row`.id = ?", id).Error
}
func (s SQLite) New(db *gorm.DB) FullTextSearch {
	return &SQLite{db: db}
}

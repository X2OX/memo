package db

import (
	"github.com/x2ox/memo/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initPostgreSQL() gorm.Dialector {
	return postgres.Open(model.Conf.DSN())
}

type PostgreSQL struct{ db *gorm.DB }

func (p PostgreSQL) Search(keywords string, offset, limit int) (arr []*model.Note, count int64) {
	p.db.Model(&model.Note{}).Where("id IN (?)", db.
		Table("note_row, to_tsquery( 'simple', ? ) query", "["+keywords+"]").
		Select("id").
		Where("note_row.tsv_content @@query").
		Order("ts_rank( note_row.tsv_content, query ) DESC").
		Offset(offset).
		Limit(limit),
	).Find(&arr)
	p.db.Model(&model.Note{}).Where("id IN (?)", p.db.
		Table("note_row, to_tsquery( 'simple', ? ) query", "["+keywords+"]").
		Select("id").
		Where("note_row.tsv_content @@query")).Count(&count)
	return
}

func (p PostgreSQL) Init() error {
	var count int64
	if err := p.db.Table("pg_class").Where("relname = ?", "note_row").Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return p.db.Exec(`CREATE TABLE "note_row" ("id" bigserial, "tsv_content" tsvector, PRIMARY KEY ("id"))`).Error
	}
	return nil
}
func (p PostgreSQL) Index() error {
	var count int64
	if err := p.db.Table("pg_class").Where("relname = ?", "tsv_content_idx").Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return p.db.Exec("CREATE INDEX tsv_content_idx ON note_row using gin(tsv_content)").Error
	}
	return nil
}
func (p PostgreSQL) ReIndex() error {
	return p.db.Exec("REINDEX INDEX tsv_content_idx;").Error
}
func (p PostgreSQL) Clean() error {
	return p.db.Exec("DROP TABLE note_row").Error
}

func (p PostgreSQL) Create(titleKeywords, contentKeywords string, id uint64) error {
	return p.db.Exec(`INSERT INTO "note_row" VALUES( ?, 
setweight( to_tsvector( 'simple', ? ), 'A' ) || 
setweight( to_tsvector( 'simple', ? ), 'B' ) )
ON CONFLICT ("id") DO UPDATE SET "tsv_content"="excluded"."tsv_content"`,
		id, titleKeywords, contentKeywords).Error
}
func (p PostgreSQL) Delete(id uint64) error {
	return p.db.Table("note_row").Delete(id).Error
}
func (p PostgreSQL) New(db *gorm.DB) FullTextSearch {
	return &PostgreSQL{db: db}
}

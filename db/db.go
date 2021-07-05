package db

import (
	"bytes"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.x2ox.com/blackdatura"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/x2ox/memo/model"
)

var (
	db *gorm.DB
)

func Init() {
	var useDialect gorm.Dialector
	if model.Conf.IsPostgreSQL() {
		useDialect = initPostgreSQL()
	} else {
		useDialect = initSQLite()
	}
	log := blackdatura.With("gorm")
	var err error
	if db, err = gorm.Open(useDialect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         blackdatura.NewGormLogger().LogMode(logger.Info),
	}); err != nil {
		log.Fatal("gorm client db fail", zap.Error(err))
	}

	if err = db.AutoMigrate(&model.Note{}, &model.Input{}); err != nil {
		log.Fatal("gorm auto migrate fail", zap.Error(err))
	}
	Search = Search.New(db)
	if err = Search.Init(); err != nil {
		log.Fatal("full text search init err", zap.Error(err))
	}
	if err = Search.Index(); err != nil {
		log.Fatal("full text search index init err", zap.Error(err))
	}
}

var (
	Note  = &noteSrv{}
	Input = &inputSrv{mux: &sync.RWMutex{}}
)

type (
	noteSrv  struct{}
	inputSrv struct {
		mux *sync.RWMutex
	}
)

func (srv *noteSrv) Find(ids []uint64) []*model.Note {
	var arr []*model.Note
	if err := db.Model(&model.Note{}).Where("id IN ?", ids).Find(&arr).Error; err != nil {
		return nil
	}
	return arr
}

func (srv *noteSrv) Query(offset, limit int) (arr []*model.Note, count int64) {
	if err := db.Model(&model.Note{}).Order("updated_at DESC").
		Offset(offset).Limit(limit).
		Find(&arr).Error; err != nil {
		return
	}

	db.Model(&model.Note{}).Count(&count)
	return
}

func (srv *noteSrv) GetWithID(id uint64) *model.Note {
	var s model.Note
	if err := db.Model(&model.Note{}).Where("id = ?", id).
		First(&s).Error; err != nil {
		return nil
	}
	return &s
}

func (srv *noteSrv) Delete(id uint64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := Search.New(tx).Delete(id); err != nil {
			return err
		}
		return tx.Delete(&model.Note{ID: id}).Error
	})
}

func (srv *noteSrv) Count() (i int64) {
	db.Model(&model.Note{}).Count(&i)
	return i
}
func (srv *noteSrv) WeekCount() (i int64) {
	db.Model(&model.Note{}).Where("created_at > ?", time.Now().Add(-7*24*time.Hour)).Count(&i)
	return i
}

func (srv *inputSrv) String() string {
	srv.mux.RLock()
	defer srv.mux.RUnlock()
	var i []string
	if err := db.Model(&model.Input{}).Select("content").Find(&i).Error; err != nil {
		return ""
	}
	return strings.Join(i, "\n")
}

func (srv *inputSrv) FindAll() []*model.Input {
	srv.mux.RLock()
	defer srv.mux.RUnlock()
	var arr []*model.Input
	if err := db.Model(&model.Input{}).Order("message_id").Find(&arr).Error; err != nil {
		return nil
	}
	return arr
}
func (srv *inputSrv) Check() bool {
	return srv.Count() != 0
}

func (srv *inputSrv) Count() (i int64) {
	srv.mux.RLock()
	defer srv.mux.RUnlock()
	if err := db.Model(&model.Input{}).Count(&i).Error; err != nil {
	}
	return i
}

func (srv *inputSrv) Submit() *model.Note {
	srv.mux.Lock()
	defer srv.mux.Unlock()

	var note *model.Note
	if err := db.Transaction(func(tx *gorm.DB) error {
		var arr []*model.Input
		if err := tx.Model(&model.Input{}).Order("message_id").Find(&arr).Error; err != nil {
			return err
		}

		var buf bytes.Buffer
		for _, v := range arr {
			buf.WriteString(v.Content)
		}

		note = model.NewNote(buf.String())
		if err := tx.Create(note).Error; err != nil {
			return err
		}
		if err := Search.New(tx).Create(note.ParticipleTitle(), note.ParticipleContent(), note.ID); err != nil {
			return err
		}

		return tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Input{}).Error
	}); err != nil {
		return nil
	}

	return note
}

func (srv *inputSrv) UpdateContent(messageID int, content string) error {
	srv.mux.Lock()
	defer srv.mux.Unlock()
	return db.Model(&model.Input{}).Where("message_id", messageID).
		Update("content", content).Error
}

func (srv *inputSrv) Add(i *model.Input) error {
	srv.mux.Lock()
	defer srv.mux.Unlock()
	return db.Create(i).Error
}

func (srv *inputSrv) Clear() error {
	srv.mux.Lock()
	defer srv.mux.Unlock()
	return db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Input{}).Error
}

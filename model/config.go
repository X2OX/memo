package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// host=127.0.0.1 user=test password=123456789 dbname=dlink port=5432 sslmode=disable TimeZone=Asia/Shanghai

type Configuration struct {
	Database        string `json:"database"`         // pg sqlite 目前仅支持，看编译的啥，数据库地址
	Domain          string `json:"domain"`           // 域名，如果不填写，默认为监听地址
	ListenAddr      string `json:"listen_addr"`      // 监听地址，默认为 :8088
	DataFolder      string `json:"data_folder"`      // 数据文件夹 默认 /data/memo/
	LogLevel        string `json:"log_level"`        // log 等级
	TelegramID      int64  `json:"telegram_id"`      // 用户的 Telegram ID
	TelegramToken   string `json:"telegram_token"`   // telegram bot token
	TelegramWebhook string `json:"telegram_webhook"` // 默认地址 /api/v1/telegram/bot/webhook

	Token struct {
		AutoUpdate uint32 `json:"auto_update"` // 自动更新 key 的时间，单位 分钟。为零不自动更新
		Preview    uint32 `json:"preview"`     // 预览的有效时间。为零不过期
		View       uint32 `json:"view"`        // 阅读的有效时间
		Share      uint32 `json:"share"`       // 分享的有效时间
	} `json:"token"`
}

var Conf Configuration

func (c Configuration) DSN() string {
	if c.IsPostgreSQL() {
		return c.Database
	}
	return filepath.Join(c.DataFolder, c.Database) + "?cache=shared&mode=rwc&_journal_mode=WAL"
}

func (c Configuration) IsSQLite() bool          { return !strings.Contains(c.Database, "host=") }
func (c Configuration) IsPostgreSQL() bool      { return strings.Contains(c.Database, "host=") }
func (c Configuration) IsWebhook() bool         { return c.TelegramWebhook != "" }
func (c Configuration) Webhook() string         { return c.Domain + c.TelegramWebhook }
func (c Configuration) TemplatesFolder() string { return filepath.Join(c.DataFolder, "/templates") }
func (c Configuration) StaticFolder() string    { return filepath.Join(c.DataFolder, "/file") }
func (c Configuration) LogFolder() string       { return filepath.Join(c.DataFolder, "/log/log") }

func mkdir(arr ...string) {
	for _, v := range arr {
		_ = os.MkdirAll(v, os.ModePerm)
	}
}

func LoadConfig() {
	var bts []byte
	if len(os.Args) > 1 {
		bts = readFile(os.Args[1])
	}
	if bts == nil {
		bts = readFile("/data/memo/config.json")
	}
	if bts == nil {
		panic("read config file error")
	}

	if err := json.Unmarshal(bts, &Conf); err != nil {
		panic(err)
	}

	if err := Conf.check(); err != nil {
		panic(err)
	}
	mkdir(Conf.DataFolder, Conf.TemplatesFolder(), Conf.StaticFolder())
	loadKey()
}

func (c *Configuration) check() error {
	if c.Domain == "" || c.TelegramID == 0 || c.TelegramToken == "" || c.TelegramWebhook == "" {
		return errors.New("config error")
	}
	if c.DataFolder == "" {
		c.DataFolder = "/data/memo/"
	}
	if c.Database == "" {
		c.Database = "db.db"
	}
	if c.ListenAddr == "" {
		c.ListenAddr = ":8088"
	}
	if c.LogLevel == "" {
		c.LogLevel = "error"
	}
	return nil
}

func readFile(s string) []byte {
	buf, err := ioutil.ReadFile(s)
	if err != nil {
		return nil
	}
	return buf
}

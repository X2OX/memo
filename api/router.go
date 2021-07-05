package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/telegram"
	"github.com/x2ox/memo/tpl"
	"go.uber.org/zap"
	"go.x2ox.com/blackdatura"
)

var log *zap.Logger

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.TrustedProxies = []string{"127.0.0.0/8", "192.168.0.0/16", "172.16.0.0/12", "10.0.0.0/8"}
	engine.NoRoute(func(c *gin.Context) { c.Status(http.StatusNotFound) })
	engine.NoMethod(func(c *gin.Context) { c.Status(http.StatusNotFound) })

	log = blackdatura.With("gin router")
	engine.Use(blackdatura.Ginzap(log))
	engine.Use(blackdatura.RecoveryWithZap(log))

	engine.SetHTMLTemplate(tpl.Load())

	engine.GET("/robots.txt", func(c *gin.Context) {
		c.String(http.StatusOK, "User-agent: *\nDisallow: /")
	})

	if model.Conf.IsWebhook() {
		engine.POST(model.Conf.TelegramWebhook, telegram.Webhook())
	}

	engine.GET("/preview/:token", previewAction)

	{
		static := engine.Group("/file")
		static.Use(authAction())
		static.Static("/", model.Conf.StaticFolder())
	}

	return engine
}

func authAction() func(c *gin.Context) {
	return func(c *gin.Context) {
		ts, err := c.Cookie("Token")
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if tk := model.ParseToken(ts); tk == nil || !tk.Valid() {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.Next()
	}
}

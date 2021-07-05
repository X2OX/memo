package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/x2ox/memo/api"
	"github.com/x2ox/memo/db"
	"github.com/x2ox/memo/model"
	"github.com/x2ox/memo/pkg/participle"
	"github.com/x2ox/memo/telegram"
	"go.uber.org/zap"
	"go.x2ox.com/blackdatura"
)

var log *zap.Logger

func init() {
	model.LoadConfig()
	blackdatura.Init(model.Conf.LogLevel, true,
		blackdatura.Lumberjack(model.Conf.LogFolder(), 1024, 30, 90, true))
	log = blackdatura.New()
	log.Info("[Memo] server starting...")
}

// TODO wait for the updated version of gin to get the fix https://github.com/gin-gonic/gin/pull/2692

func main() {
	participle.Init(model.Conf.DataFolder)
	db.Init()
	telegram.Init()

	router := api.Router()
	// server := &http.Server{
	// 	Addr:    model.Conf.ListenAddr,
	// 	Handler: router,
	// }
	// handleSignal(server)
	//
	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatal("[Memo] server listen error", zap.Error(err))
	// }
	if err := router.Run(model.Conf.ListenAddr); err != nil {
		log.Fatal("[Memo] server listen error", zap.Error(err))
	}
}

// handleSignal handles system signal for graceful shutdown.
func handleSignal(server *http.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		sig := <-c
		log.Info("[Memo] server exit signal", zap.Any("signal notify", sig))

		_ = server.Close()
		telegram.Close()

		log.Info("[Memo] server exit", zap.Time("time", time.Now()))
		os.Exit(0)
	}()
}

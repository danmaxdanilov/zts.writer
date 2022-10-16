package main

import (
	"flag"
	"log"

	"github.com/danmaxdanilov/zts.shared/pkg/logger"
	"github.com/danmaxdanilov/zts.writer/config"
	"github.com/danmaxdanilov/zts.writer/internal/server"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("zts.writer")

	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}

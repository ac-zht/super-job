package main

import (
	"github.com/ac-zht/super-job/admin/internal/service"
	"go.uber.org/zap"
)

func main() {
	initLogger()
	service.InitEnv()

	web := InitWeb()
	web.Run(":9100")
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

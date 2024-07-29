package main

import (
	"context"
	"fmt"
	"github.com/ac-zht/super-job/scheduler/ioc"
	"go.uber.org/zap"
)

func main() {
	initLogger()
	scheduler := InitScheduler()
	ctx := context.Background()
	go func() {
		err := scheduler.Start(ctx)
		if err != nil {
			zap.L().Error(fmt.Sprintf("自动调度错误终止#%v", err))
		}
	}()
	web := ioc.InitWebServer(scheduler)
	err := web.Run("9200")
	if err != nil {
		zap.L().Error(fmt.Sprintf("web服务执行错误#%v", err))
	}
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

//go:build wireinject

package main

import (
	"github.com/ac-zht/super-job/admin/internal/repository"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/internal/web"
	"github.com/ac-zht/super-job/admin/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWeb() *gin.Engine {
	wire.Build(
		ioc.InitDB,

		dao.NewExecutorDAO,
		dao.NewTaskDAO,
		dao.NewSettingDAO,

		repository.NewExecutorRepository,
		repository.NewTaskRepository,
		repository.NewSettingRepository,

		service.NewExecutorService,
		service.NewTaskService,
		service.NewSettingService,
		service.NewInstallService,

		web.NewExecutorHandler,
		web.NewTaskHandler,
		web.NewSettingHandler,
		web.NewInstallHandler,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}

//go:build wireinject

package main

import (
	"github.com/ac-zht/gotools/option"
	"github.com/ac-zht/super-job/admin/internal/repository"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/internal/web"
	"github.com/ac-zht/super-job/admin/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebSettingService() service.WebSettingService {
	wire.Build(
		commonProvider,
		service.NewWebSettingService,
	)
	return &service.WebSettingSvc{}
}

var commonProvider = wire.NewSet(
	NewBaseModelOption,
	dao.NewBaseModel,

	dao.NewSettingDAO,
	dao.NewUserDAO,
	dao.NewInstallDAO,

	repository.NewSettingRepository,
	repository.NewUserRepository,
	repository.NewInstallRepository,

	repository.NewWebSettingRepository,

	service.NewInstallService,
)

func InitInstallService() service.InstallService {
	wire.Build(commonProvider)
	return &service.InstallSvc{}
}

func NewBaseModelOption() []option.Option[dao.BaseDbModel] {
	return []option.Option[dao.BaseDbModel]{func(m *dao.BaseDbModel) {}}
}

func InitWeb() *gin.Engine {
	wire.Build(
		//ioc.InitDB,
		commonProvider,
		dao.NewExecutorDAO,

		dao.NewTaskDAO,
		repository.NewExecutorRepository,

		repository.NewTaskRepository,
		service.NewExecutorService,
		service.NewTaskService,
		service.NewSettingService,

		web.NewExecutorHandler,
		web.NewTaskHandler,
		web.NewSettingHandler,
		web.NewInstallHandler,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}

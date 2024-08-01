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

func InitWebSettingService(repo repository.WebSettingRepository) service.WebSettingService {
	wire.Build(
		service.NewWebSettingService,
	)
	return &service.WebSettingSvc{}
}

func InitInstallService(base dao.BaseModel,
	setRepo repository.SettingRepository,
	webSettingRepo repository.WebSettingRepository,
	userRepo repository.UserRepository) service.InstallService {
	wire.Build(
		dao.NewInstallDAO,
		repository.NewInstallRepository,
		service.NewInstallService,
	)
	return &service.InstallSvc{}
}

func NewBaseModelOption() []option.Option[dao.BaseDbModel] {
	return []option.Option[dao.BaseDbModel]{func(m *dao.BaseDbModel) {}}
}

func InitWeb() *gin.Engine {
	wire.Build(
		//ioc.InitDB,
		NewBaseModelOption,
		dao.NewBaseModel,

		dao.NewExecutorDAO,
		dao.NewTaskDAO,
		dao.NewSettingDAO,
		dao.NewUserDAO,

		repository.NewExecutorRepository,
		repository.NewTaskRepository,
		repository.NewSettingRepository,
		repository.NewUserRepository,
		repository.NewWebSettingRepository,

		service.NewExecutorService,
		service.NewTaskService,
		service.NewSettingService,
		//InitWebSettingService,
		InitInstallService,

		web.NewExecutorHandler,
		web.NewTaskHandler,
		web.NewSettingHandler,
		web.NewInstallHandler,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}

package main

import (
	"github.com/ac-zht/super-job/scheduler/internal/repository"
	"github.com/ac-zht/super-job/scheduler/internal/repository/dao"
	"github.com/ac-zht/super-job/scheduler/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitScheduler() *gin.Engine {
	wire.Build(
		ioc.InitDB,

		dao.NewSettingDAO,
		dao.NewTaskDAO,
		dao.NewTaskLogDAO,

		repository.NewSettingRepository,
		repository.NewSettingRepository,
		repository.NewSettingRepository,
	)
	return new(gin.Engine)
}

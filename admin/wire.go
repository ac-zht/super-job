//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/zc-zht/super-job/admin/internal/repository"
	"github.com/zc-zht/super-job/admin/internal/repository/dao"
	"github.com/zc-zht/super-job/admin/internal/service"
	"github.com/zc-zht/super-job/admin/internal/web"
	"github.com/zc-zht/super-job/admin/ioc"
)

func InitWeb() *gin.Engine {
	wire.Build(
		ioc.InitDB,

		dao.NewExecutorDAO,
		dao.NewJobDAO,

		repository.NewExecutorRepository,
		repository.NewJobRepository,

		service.NewExecutorService,
		service.NewJobService,

		web.NewExecutorHandler,
		web.NewJobHandler,

		ioc.InitWebServer,
	)
	return new(gin.Engine)
}

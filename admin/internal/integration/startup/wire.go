//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/zc-zht/super-job/admin/internal/repository"
	"github.com/zc-zht/super-job/admin/internal/repository/dao"
	"github.com/zc-zht/super-job/admin/internal/service"
	"github.com/zc-zht/super-job/admin/internal/web"
)

var thirdProvider = wire.NewSet(
	InitTestDB,
)

var executorSvcProvider = wire.NewSet(
	dao.NewExecutorDAO,
	repository.NewExecutorRepository,
	service.NewExecutorService,
	web.NewExecutorHandler)

var jobSvcProvider = wire.NewSet(
	dao.NewJobDAO,
	repository.NewJobRepository,
	service.NewJobService,
	web.NewJobHandler)

func InitWeb() *gin.Engine {
	wire.Build(
		thirdProvider,
		executorSvcProvider,
		jobSvcProvider,
		InitTestWebServer,
	)
	return new(gin.Engine)
}

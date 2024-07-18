//go:build wireinject

package startup

import (
	"github.com/ac-zht/super-job/admin/internal/repository"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(
	InitTestDB,
)

var executorSvcProvider = wire.NewSet(
	dao.NewExecutorDAO,
	repository.NewExecutorRepository,
	service.NewExecutorService,
	web.NewExecutorHandler)

var taskSvcProvider = wire.NewSet(
	dao.NewTaskDAO,
	repository.NewTaskRepository,
	service.NewTaskService,
	web.NewTaskHandler)

func InitWeb() *gin.Engine {
	wire.Build(
		thirdProvider,
		executorSvcProvider,
		taskSvcProvider,
		InitTestWebServer,
	)
	return new(gin.Engine)
}

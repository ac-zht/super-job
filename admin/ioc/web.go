package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/zc-zht/super-job/admin/internal/web"
)

func InitWebServer(jobHandler *web.JobHandler,
	executorHandler *web.ExecutorHandler,
	settingHandler *web.SettingHandler,
	installHandler *web.InstallHandler) *gin.Engine {
	server := gin.Default()
	jobHandler.RegisterRoutes(server)
	executorHandler.RegisterRoutes(server)
	settingHandler.RegisterRoutes(server)
	installHandler.RegisterRoutes(server)
	return server
}

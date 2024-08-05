package ioc

import (
	"github.com/ac-zht/super-job/admin/internal/web"
	"github.com/gin-gonic/gin"
)

func InitWebServer(taskHandler *web.TaskHandler,
	executorHandler *web.ExecutorHandler,
	settingHandler *web.SettingHandler,
	installHandler *web.InstallHandler) *gin.Engine {
	server := gin.Default()
	taskHandler.RegisterRoutes(server)
	executorHandler.RegisterRoutes(server)
	settingHandler.RegisterRoutes(server)
	installHandler.RegisterRoutes(server)
	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

func userAuthHdl() {

}

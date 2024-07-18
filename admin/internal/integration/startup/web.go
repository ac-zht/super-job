package startup

import (
	"github.com/ac-zht/super-job/admin/internal/web"
	"github.com/gin-gonic/gin"
)

func InitTestWebServer(taskHandler *web.TaskHandler,
	executorHandler *web.ExecutorHandler) *gin.Engine {
	server := gin.Default()
	taskHandler.RegisterRoutes(server)
	executorHandler.RegisterRoutes(server)
	return server
}

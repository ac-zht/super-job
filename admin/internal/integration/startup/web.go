package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/zc-zht/super-job/admin/internal/web"
)

func InitTestWebServer(jobHandler *web.JobHandler,
	executorHandler *web.ExecutorHandler) *gin.Engine {
	server := gin.Default()
	jobHandler.RegisterRoutes(server)
	executorHandler.RegisterRoutes(server)
	return server
}

package ioc

import (
	"github.com/ac-zht/super-job/scheduler/internal/web"
	"github.com/gin-gonic/gin"
)

func InitWebServer(scheduler *web.Scheduler) *gin.Engine {
	server := gin.Default()
	scheduler.RegisterRoutes(server)
	return server
}

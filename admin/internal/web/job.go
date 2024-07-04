package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zht-account/gotools/pool"
	"github.com/zht-account/super-job/admin/internal/service"
	"github.com/zht-account/super-job/admin/pkg/logger"
	"time"
)

type JobHandler struct {
	svc          service.JobService
	failInterval time.Duration
	dbTimeout    time.Duration
	quickPool    *pool.OnDemandBlockTaskPool
	l            logger.Logger
}

func NewJobHandler(svc service.JobService) *JobHandler {
	return &JobHandler{
		svc: svc,
	}
}

func (h *JobHandler) Store() {
	return
}

func (h *JobHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/job")
	ug.POST("/store", h.Store)
}

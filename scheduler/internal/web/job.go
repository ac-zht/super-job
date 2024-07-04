package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zht-account/gotools/pool"
	"github.com/zht-account/super-job/scheduler/internal/domain"
	"github.com/zht-account/super-job/scheduler/internal/errs"
	"github.com/zht-account/super-job/scheduler/internal/service"
	"github.com/zht-account/super-job/scheduler/pkg/ginx"
	"github.com/zht-account/super-job/scheduler/pkg/logger"
	"net/http"
	"time"
)

type Scheduler struct {
	svc          service.JobService
	failInterval time.Duration
	dbTimeout    time.Duration
	quickPool    *pool.OnDemandBlockTaskPool
	l            logger.Logger
}

func NewScheduler(svc service.JobService, interval, dt time.Duration, qp *pool.OnDemandBlockTaskPool, l logger.Logger) *Scheduler {
	return &Scheduler{
		svc:          svc,
		failInterval: interval,
		dbTimeout:    dt,
		quickPool:    qp,
		l:            l,
	}
}

func (h *Scheduler) Start(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		dbCtx, cancel := context.WithTimeout(ctx, h.dbTimeout)
		j, err := h.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			time.Sleep(h.failInterval)
			continue
		}
		fmt.Println(j)
		//放入到线程池
	}
}

func (h *Scheduler) RegisterJob(ctx *gin.Context) {
	type Req struct {
		Name       string `json:"name"`
		Executor   string `json:"executor"`
		Cfg        string `json:"cfg"`
		Expression string `json:"expression"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := h.svc.AddJob(ctx, domain.Job{
		Name:       req.Name,
		Executor:   req.Executor,
		Cfg:        req.Cfg,
		Expression: req.Expression,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.JobInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "新增成功",
	})
	return
}

func (h *Scheduler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/job")
	ug.POST("/register", h.RegisterJob)
}

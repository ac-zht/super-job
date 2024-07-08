package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zc-zht/gotools/pool"
	"github.com/zc-zht/super-job/scheduler/internal/service"
	"github.com/zc-zht/super-job/scheduler/pkg/logger"
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

func (h *Scheduler) Start(ctx *gin.Context) error {
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

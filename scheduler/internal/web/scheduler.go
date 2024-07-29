package web

import (
	"context"
	"github.com/ac-zht/gotools/option"
	"github.com/ac-zht/gotools/pool"
	"github.com/ac-zht/super-job/scheduler/internal/service"
	"github.com/gin-gonic/gin"
	"time"
)

type Scheduler struct {
	svc          service.JobService
	failInterval time.Duration
	dbTimeout    time.Duration
	quickPool    *pool.OnDemandBlockTaskPool
}

func NewScheduler(svc service.JobService, qp *pool.OnDemandBlockTaskPool, opts ...option.Option[Scheduler]) *Scheduler {
	scheduler := &Scheduler{
		svc:          svc,
		failInterval: time.Second,
		dbTimeout:    time.Second,
		quickPool:    qp,
	}
	option.Apply[Scheduler](scheduler, opts...)
	return scheduler
}

func WithFailInterval(interval time.Duration) option.Option[Scheduler] {
	return func(s *Scheduler) {
		s.failInterval = interval
	}
}

func WithDbTimeout(dt time.Duration) option.Option[Scheduler] {
	return func(s *Scheduler) {
		s.dbTimeout = dt
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
		//放入到线程池
		job := h.svc.CreateJob(j)
		err = h.quickPool.Submit(ctx, pool.TaskFunc(job))
		//阻塞自旋
		if err != nil {
			time.Sleep(time.Second * 5)
			continue
		}
	}
}

func (h *Scheduler) ExecTask(ctx *gin.Context) {

}

func (h *Scheduler) Stop(ctx *gin.Context) {

}

func (h *Scheduler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/api/dispatch")
	ug.GET("/task/:id", h.ExecTask)
	ug.GET("/stop", h.Stop)
}

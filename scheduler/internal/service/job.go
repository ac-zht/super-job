package service

import (
	"context"
	"github.com/zht-account/super-job/scheduler/internal/domain"
	"github.com/zht-account/super-job/scheduler/internal/repository"
	"github.com/zht-account/super-job/scheduler/pkg/logger"
	"time"
)

//go:generate mockgen -source=./cron_job.go -package=svcmocks -destination=mocks/cron_job.mocks.go JobService
type JobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, job domain.Job) error
	AddJob(ctx context.Context, job domain.Job) error
}

type cronJobService struct {
	repo            repository.JobRepository
	l               logger.Logger
	refreshInterval time.Duration
}

func NenJobService(
	repo repository.JobRepository,
	l logger.Logger) JobService {
	return &cronJobService{
		repo:            repo,
		l:               l,
		refreshInterval: time.Second * 10,
	}
}

func (c *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	//ch := make(chan struct{})
	//go func() {
	//	ticker := time.NewTicker(c.refreshInterval)
	//	defer ticker.Stop()
	//	for {
	//		select {
	//		case <-ch:
	//			return
	//		case <-ticker.C:
	//			c.refresh(j.Id)
	//		}
	//	}
	//}()
	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		for range ticker.C {
			c.refresh(j.Id)
		}
	}()

	j.CancelFunc = func() {
		//close(ch)
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.repo.Release(ctx, j.Id)
		if err != nil {
			c.l.Error("释放任务失败",
				logger.Error(err),
				logger.Int64("id", j.Id))
		}
	}
	return j, nil
}

func (c *cronJobService) ResetNextTime(ctx context.Context, job domain.Job) error {
	t := job.Next(time.Now())
	if !t.IsZero() {
		return c.repo.UpdateNextTime(ctx, job.Id, t)
	}
	return nil
}

func (c *cronJobService) AddJob(ctx context.Context, job domain.Job) error {
	job.NextTime = job.Next(time.Now())
	return c.repo.AddJob(ctx, job)
}

func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.repo.UpdateUtime(ctx, id)
	if err != nil {
		c.l.Error("续约失败",
			logger.Int64("jid", id),
			logger.Error(err))
	}
}

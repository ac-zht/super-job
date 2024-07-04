package service

import (
	"context"
	"github.com/zht-account/super-job/admin/internal/domain"
	"github.com/zht-account/super-job/admin/internal/repository"
	"github.com/zht-account/super-job/admin/pkg/logger"
	"time"
)

//go:generate mockgen -source=./cron_job.go -package=svcmocks -destination=mocks/cron_job.mocks.go JobService
type JobService interface {
	List(ctx context.Context, offset, limit int) ([]domain.Job, error)
	Store(ctx context.Context, j domain.Job) error
	Delete(ctx context.Context, id int) error
}

type jobService struct {
	repo repository.JobRepository
	l    logger.Logger
}

func (c *jobService) List(ctx context.Context, offset, limit int) ([]domain.Job, error) {
	return c.repo.List(ctx, offset, limit)
}

func (c *jobService) Store(ctx context.Context, j domain.Job) error {
	j.NextTime = j.Next(time.Now())
	return c.repo.Store(ctx, j)
}

func (c *jobService) Delete(ctx context.Context, id int) error {
	return c.repo.Delete(ctx, id)
}

func NewJobService(repo repository.JobRepository, l logger.Logger) JobService {
	return &jobService{
		repo: repo,
		l:    l,
	}
}

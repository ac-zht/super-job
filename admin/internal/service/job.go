package service

import (
	"context"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/repository"
	"time"
)

//go:generate mockgen -source=./cron_job.go -package=svcmocks -destination=mocks/cron_job.mocks.go JobService
type JobService interface {
	List(ctx context.Context, offset, limit int) ([]domain.Job, error)
	Save(ctx context.Context, j domain.Job) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type jobService struct {
	repo repository.JobRepository
	//l    logger.Logger
}

func NewJobService(repo repository.JobRepository) JobService {
	return &jobService{
		repo: repo,
	}
}

func (svc *jobService) List(ctx context.Context, offset, limit int) ([]domain.Job, error) {
	return svc.repo.List(ctx, offset, limit)
}

func (svc *jobService) Save(ctx context.Context, j domain.Job) (int64, error) {
	j.NextTime = j.Next(time.Now())
	if j.Id > 0 {
		err := svc.update(ctx, j)
		return j.Id, err
	}
	return svc.create(ctx, j)
}

func (svc *jobService) create(ctx context.Context, j domain.Job) (int64, error) {
	return svc.repo.Create(ctx, j)
}

func (svc *jobService) update(ctx context.Context, j domain.Job) error {
	return svc.repo.Update(ctx, j)
}

func (svc *jobService) Delete(ctx context.Context, id int64) error {
	return svc.repo.Delete(ctx, id)
}

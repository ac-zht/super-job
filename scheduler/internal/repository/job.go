package repository

import (
	"context"
	"github.com/zc-zht/super-job/scheduler/internal/domain"
	"github.com/zc-zht/super-job/scheduler/internal/repository/dao"
	"time"
)

var ErrNoMoreJob = dao.ErrNoMoreJob

//go:generate mockgen -source=./job.go -package=repomocks -destination=mocks/job.mock.go JobRepository
type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	AddJob(ctx context.Context, j domain.Job) error
}

type PreemptJobRepository struct {
	dao dao.JobDAO
}

func (p *PreemptJobRepository) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return p.toDomain(j), nil
}

func (p *PreemptJobRepository) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	return p.dao.UpdateNextTime(ctx, id, t)
}

func (p *PreemptJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return p.dao.UpdateUtime(ctx, id)
}

func (p *PreemptJobRepository) Release(ctx context.Context, id int64) error {
	return p.dao.Release(ctx, id)
}

func (p *PreemptJobRepository) AddJob(ctx context.Context, j domain.Job) error {
	return p.dao.Insert(ctx, p.toEntity(j))
}

func (p *PreemptJobRepository) toEntity(j domain.Job) dao.Job {
	return dao.Job{
		Id:         j.Id,
		Name:       j.Name,
		Executor:   j.Executor,
		Cfg:        j.Cfg,
		Expression: j.Expression,
		NextTime:   j.NextTime.UnixMilli(),
	}
}

func (p *PreemptJobRepository) toDomain(j dao.Job) domain.Job {
	return domain.Job{
		Id:         j.Id,
		Name:       j.Name,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		Executor:   j.Executor,
		NextTime:   time.UnixMilli(j.NextTime),
	}
}

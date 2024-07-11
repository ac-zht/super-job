package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/repository/dao"
	"time"
)

var ErrNoMoreJob = dao.ErrNoMoreJob

//go:generate mockgen -source=./job.go -package=repomocks -destination=mocks/job.mock.go JobRepository
type JobRepository interface {
	List(ctx context.Context, offset, limit int) ([]domain.Job, error)
	GetById(ctx context.Context, id int64) (domain.Job, error)
	Create(ctx context.Context, j domain.Job) (int64, error)
	Update(ctx context.Context, job domain.Job) error
	Delete(ctx context.Context, id int64) error
}

type PreemptJobRepository struct {
	dao dao.JobDAO
}

func NewJobRepository(dao dao.JobDAO) JobRepository {
	return &PreemptJobRepository{
		dao: dao,
	}
}

func (p *PreemptJobRepository) List(ctx context.Context, offset, limit int) ([]domain.Job, error) {
	jobs, err := p.dao.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Job, domain.Job](jobs, func(idx int, src dao.Job) domain.Job {
		return p.toDomain(src)
	}), nil
}

func (p *PreemptJobRepository) Create(ctx context.Context, j domain.Job) (int64, error) {
	return p.dao.Insert(ctx, p.toEntity(j))
}

func (p *PreemptJobRepository) Update(ctx context.Context, job domain.Job) error {
	return p.dao.Update(ctx, p.toEntity(job))
}

func (p *PreemptJobRepository) Delete(ctx context.Context, id int64) error {
	return p.dao.Delete(ctx, id)
}

func (p *PreemptJobRepository) GetById(ctx context.Context, id int64) (domain.Job, error) {
	job, err := p.dao.GetById(ctx, id)
	return p.toDomain(job), err
}

func (p *PreemptJobRepository) toEntity(j domain.Job) dao.Job {
	return dao.Job{
		Id:            j.Id,
		ExecId:        j.ExecId,
		Name:          j.Name,
		Expression:    j.Expression,
		Cfg:           j.Cfg,
		Status:        j.Status,
		NextTime:      j.NextTime.UnixMilli(),
		Protocol:      j.Protocol.ToUint8(),
		HttpMethod:    j.HttpMethod.ToUint8(),
		Multi:         j.Multi,
		Timeout:       j.Timeout,
		RetryTimes:    j.RetryTimes,
		RetryInterval: j.RetryInterval,
		Creator:       j.Creator,
		Updater:       j.Updater,
	}
}

func (p *PreemptJobRepository) toDomain(j dao.Job) domain.Job {
	executor := &executorRepository{}
	return domain.Job{
		Id:            j.Id,
		ExecId:        j.ExecId,
		Name:          j.Name,
		Expression:    j.Expression,
		Cfg:           j.Cfg,
		Status:        j.Status,
		NextTime:      time.UnixMilli(j.NextTime),
		Executor:      executor.toDomain(j.Executor),
		Protocol:      domain.JobProtocol(j.Protocol),
		HttpMethod:    domain.HttpMethod(j.HttpMethod),
		Multi:         j.Multi,
		Timeout:       j.Timeout,
		RetryTimes:    j.RetryTimes,
		RetryInterval: j.RetryInterval,
		Creator:       j.Creator,
		Updater:       j.Updater,
		Ctime:         j.Ctime,
		Utime:         j.Utime,
	}
}

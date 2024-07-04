package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/zht-account/super-job/admin/internal/domain"
	"github.com/zht-account/super-job/admin/internal/repository/dao"
	"time"
)

var ErrNoMoreJob = dao.ErrNoMoreJob

//go:generate mockgen -source=./job.go -package=repomocks -destination=mocks/job.mock.go JobRepository
type JobRepository interface {
	List(ctx context.Context, offset, limit int) ([]domain.Job, error)
	Store(ctx context.Context, j domain.Job) error
	Delete(ctx context.Context, id int) error
}

type PreemptJobRepository struct {
	dao dao.JobDAO
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

func (p *PreemptJobRepository) Store(ctx context.Context, j domain.Job) error {
	return p.dao.Insert(ctx, p.toEntity(j))
}

func (p *PreemptJobRepository) Delete(ctx context.Context, id int) error {
	return p.dao.Delete(ctx, id)
}

func (p *PreemptJobRepository) toEntity(j domain.Job) dao.Job {
	return dao.Job{
		Id:         j.Id,
		ExecId:     j.ExecId,
		Name:       j.Name,
		Protocol:   j.Protocol.ToUint8(),
		Expression: j.Expression,
		Cfg:        j.Cfg,
		NextTime:   j.NextTime.UnixMilli(),
	}
}

func (p *PreemptJobRepository) toDomain(j dao.Job) domain.Job {
	return domain.Job{
		Id:         j.Id,
		ExecId:     j.ExecId,
		Name:       j.Name,
		Protocol:   domain.JobProtocol(j.Protocol),
		Expression: j.Expression,
		Cfg:        j.Cfg,
		NextTime:   time.UnixMilli(j.NextTime),
	}
}

package repository

import (
	"context"
	"github.com/ac-zht/super-job/scheduler/internal/domain"
	"github.com/ac-zht/super-job/scheduler/internal/repository/dao"
	"time"
)

var ErrNoMoreTask = dao.ErrNoMoreTask

//go:generate mockgen -source=./task.go -package=repomocks -destination=mocks/task.mock.go TaskRepository
type TaskRepository interface {
	Preempt(ctx context.Context) (domain.Task, error)
	UpdateNextTime(ctx context.Context, id int64, t time.Time) error
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	AddTask(ctx context.Context, j domain.Task) error
}

type PreemptTaskRepository struct {
	dao dao.TaskDAO
}

func (p *PreemptTaskRepository) Preempt(ctx context.Context) (domain.Task, error) {
	j, err := p.dao.Preempt(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	return p.toDomain(j), nil
}

func (p *PreemptTaskRepository) UpdateNextTime(ctx context.Context, id int64, t time.Time) error {
	return p.dao.UpdateNextTime(ctx, id, t)
}

func (p *PreemptTaskRepository) UpdateUtime(ctx context.Context, id int64) error {
	return p.dao.UpdateUtime(ctx, id)
}

func (p *PreemptTaskRepository) Release(ctx context.Context, id int64) error {
	return p.dao.Release(ctx, id)
}

func (p *PreemptTaskRepository) toEntity(j domain.Task) dao.Task {
	return dao.Task{
		Id:         j.Id,
		Name:       j.Name,
		Cfg:        j.Cfg,
		Expression: j.Expression,
		NextTime:   j.NextTime.UnixMilli(),
	}
}

func (p *PreemptTaskRepository) toDomain(j dao.Task) domain.Task {
	return domain.Task{
		Id:         j.Id,
		Name:       j.Name,
		Expression: j.Expression,
		Cfg:        j.Cfg,
		NextTime:   time.UnixMilli(j.NextTime),
	}
}

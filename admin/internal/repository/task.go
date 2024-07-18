package repository

import (
	"context"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

var ErrNoMoreTask = dao.ErrNoMoreTask

//go:generate mockgen -source=./task.go -package=repomocks -destination=mocks/task.mock.go TaskRepository
type TaskRepository interface {
	List(ctx context.Context, offset, limit int) ([]domain.Task, error)
	GetById(ctx context.Context, id int64) (domain.Task, error)
	Create(ctx context.Context, j domain.Task) (int64, error)
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id int64) error
}

type PreemptTaskRepository struct {
	dao dao.TaskDAO
}

func NewTaskRepository(dao dao.TaskDAO) TaskRepository {
	return &PreemptTaskRepository{
		dao: dao,
	}
}

func (p *PreemptTaskRepository) List(ctx context.Context, offset, limit int) ([]domain.Task, error) {
	tasks, err := p.dao.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Task, domain.Task](tasks, func(idx int, src dao.Task) domain.Task {
		return p.toDomain(src)
	}), nil
}

func (p *PreemptTaskRepository) Create(ctx context.Context, j domain.Task) (int64, error) {
	return p.dao.Insert(ctx, p.toEntity(j))
}

func (p *PreemptTaskRepository) Update(ctx context.Context, task domain.Task) error {
	return p.dao.Update(ctx, p.toEntity(task))
}

func (p *PreemptTaskRepository) Delete(ctx context.Context, id int64) error {
	return p.dao.Delete(ctx, id)
}

func (p *PreemptTaskRepository) GetById(ctx context.Context, id int64) (domain.Task, error) {
	task, err := p.dao.GetById(ctx, id)
	return p.toDomain(task), err
}

func (p *PreemptTaskRepository) toEntity(j domain.Task) dao.Task {
	return dao.Task{
		Id:               j.Id,
		ExecId:           j.ExecId,
		Name:             j.Name,
		Expression:       j.Expression,
		Cfg:              j.Cfg,
		Status:           j.Status,
		NextTime:         j.NextTime.UnixMilli(),
		Multi:            j.Multi,
		Protocol:         j.Protocol.ToUint8(),
		HttpMethod:       j.HttpMethod.ToUint8(),
		ExecutorHandler:  j.ExecutorHandler,
		Command:          j.Command,
		Timeout:          j.Timeout,
		RetryTimes:       j.RetryTimes,
		RetryInterval:    j.RetryInterval,
		NotifyStatus:     j.NotifyStatus.ToUint8(),
		NotifyType:       j.NotifyType.ToUint8(),
		NotifyReceiverId: j.NotifyReceiverId,
		NotifyKeyword:    j.NotifyKeyword,
		Creator:          j.Creator,
		Updater:          j.Updater,
	}
}

func (p *PreemptTaskRepository) toDomain(j dao.Task) domain.Task {
	executor := &executorRepository{}
	return domain.Task{
		Id:               j.Id,
		ExecId:           j.ExecId,
		Name:             j.Name,
		Expression:       j.Expression,
		Cfg:              j.Cfg,
		Status:           j.Status,
		NextTime:         time.UnixMilli(j.NextTime),
		Multi:            j.Multi,
		Protocol:         domain.TaskProtocol(j.Protocol),
		HttpMethod:       domain.HttpMethod(j.HttpMethod),
		ExecutorHandler:  j.ExecutorHandler,
		Command:          j.Command,
		Timeout:          j.Timeout,
		RetryTimes:       j.RetryTimes,
		RetryInterval:    j.RetryInterval,
		NotifyStatus:     domain.NotifyStatus(j.NotifyStatus),
		NotifyType:       domain.NotifyType(j.NotifyType),
		NotifyReceiverId: j.NotifyReceiverId,
		NotifyKeyword:    j.NotifyKeyword,
		Creator:          j.Creator,
		Updater:          j.Updater,
		Ctime:            j.Ctime,
		Utime:            j.Utime,
		Executor:         executor.toDomain(j.Executor),
	}
}

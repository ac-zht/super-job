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
}

type PreemptTaskRepository struct {
	dao dao.TaskDAO
}

func NewTaskRepository(dao dao.TaskDAO) TaskRepository {
	return &PreemptTaskRepository{
		dao: dao,
	}
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
		Id:                    j.Id,
		ExecId:                j.ExecId,
		Name:                  j.Name,
		Cfg:                   j.Cfg,
		Expression:            j.Expression,
		NextTime:              j.NextTime.UnixMilli(),
		Status:                j.Status,
		Multi:                 j.Multi,
		Protocol:              j.Protocol.ToUint8(),
		HttpMethod:            j.HttpMethod.ToUint8(),
		ExecutorHandler:       j.ExecutorHandler,
		Command:               j.Command,
		ExecutorRouteStrategy: j.ExecutorRouteStrategy,
		Timeout:               int64(j.Timeout),
		RetryTimes:            j.RetryTimes,
		RetryInterval:         j.RetryInterval,
		NotifyStatus:          j.NotifyStatus.ToUint8(),
		NotifyType:            j.NotifyType.ToUint8(),
		NotifyReceiverId:      j.NotifyReceiverId,
		NotifyKeyword:         j.NotifyKeyword,
		Utime:                 j.Utime,
		Executor:              domain.ToEntity(j.Executor),
	}
}

func (p *PreemptTaskRepository) toDomain(j dao.Task) domain.Task {
	return domain.Task{
		Id:                    j.Id,
		ExecId:                j.ExecId,
		Name:                  j.Name,
		Cfg:                   j.Cfg,
		Expression:            j.Expression,
		NextTime:              time.UnixMilli(j.NextTime),
		Status:                j.Status,
		Multi:                 j.Multi,
		Protocol:              domain.TaskProtocol(j.Protocol),
		HttpMethod:            domain.HttpMethod(j.HttpMethod),
		ExecutorHandler:       j.ExecutorHandler,
		Command:               j.Command,
		ExecutorRouteStrategy: j.ExecutorRouteStrategy,
		Timeout:               time.Duration(j.Timeout),
		RetryTimes:            j.RetryTimes,
		RetryInterval:         j.RetryInterval,
		NotifyStatus:          domain.NotifyStatus(j.NotifyStatus),
		NotifyType:            domain.NotifyType(j.NotifyType),
		NotifyReceiverId:      j.NotifyReceiverId,
		NotifyKeyword:         j.NotifyKeyword,
		Utime:                 j.Utime,
		Executor:              domain.ToDomain(j.Executor),
	}
}

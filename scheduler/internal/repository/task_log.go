package repository

import (
	"context"
	"github.com/ac-zht/super-job/scheduler/internal/domain"
	"github.com/ac-zht/super-job/scheduler/internal/repository/dao"
)

type TaskLogRepository interface {
	Create(ctx context.Context, taskLog domain.TaskLog) (int64, error)
	UpdateById(ctx context.Context, id int64, data domain.CommonMap) error
}

type taskLogRepository struct {
	dao dao.TaskLogDAO
}

func NewTaskLogRepository(dao dao.TaskLogDAO) TaskLogRepository {
	return &taskLogRepository{
		dao: dao,
	}
}

func (repo *taskLogRepository) Create(ctx context.Context, taskLog domain.TaskLog) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(taskLog))
}

func (repo *taskLogRepository) UpdateById(ctx context.Context, id int64, data domain.CommonMap) error {
	return repo.dao.UpdateById(ctx, id, dao.CommonMap(data))
}

func (repo *taskLogRepository) toEntity(taskLog domain.TaskLog) dao.TaskLog {
	return dao.TaskLog{
		TaskId:        taskLog.Id,
		ExecId:        taskLog.ExecId,
		Name:          taskLog.Name,
		Spec:          taskLog.Spec,
		SchedulerAddr: taskLog.SchedulerAddr,
		Protocol:      taskLog.Protocol,
		Command:       taskLog.Command,
		ExecutorMsg:   "",
		Timeout:       int64(taskLog.Timeout),
		RetryTimes:    taskLog.RetryTimes,
		StartTime:     taskLog.StartTime.UnixMilli(),
		EndTime:       taskLog.EndTime.UnixMilli(),
		Status:        taskLog.Status,
		TotalTime:     taskLog.TotalTime,
		Result:        taskLog.Result,
		NotifyStatus:  taskLog.NotifyStatus,
	}
}

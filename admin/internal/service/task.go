package service

import (
	"context"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository"
	"time"
)

//go:generate mockgen -source=./cron_task.go -package=svcmocks -destination=mocks/cron_task.mocks.go TaskService
type TaskService interface {
	List(ctx context.Context, offset, limit int) ([]domain.Task, error)
	GetById(ctx context.Context, id int64) (domain.Task, error)
	Save(ctx context.Context, j domain.Task) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{
		repo: repo,
	}
}

func (svc *taskService) List(ctx context.Context, offset, limit int) ([]domain.Task, error) {
	return svc.repo.List(ctx, offset, limit)
}

func (svc *taskService) GetById(ctx context.Context, id int64) (domain.Task, error) {
	return svc.repo.GetById(ctx, id)
}

func (svc *taskService) Save(ctx context.Context, j domain.Task) (int64, error) {
	j.NextTime = j.Next(time.Now())
	if j.Id > 0 {
		err := svc.update(ctx, j)
		return j.Id, err
	}
	return svc.create(ctx, j)
}

func (svc *taskService) create(ctx context.Context, j domain.Task) (int64, error) {
	return svc.repo.Create(ctx, j)
}

func (svc *taskService) update(ctx context.Context, j domain.Task) error {
	return svc.repo.Update(ctx, j)
}

func (svc *taskService) Delete(ctx context.Context, id int64) error {
	return svc.repo.Delete(ctx, id)
}

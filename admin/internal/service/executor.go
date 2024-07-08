package service

import (
	"context"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/repository"
)

type ExecutorService interface {
	List(ctx context.Context, offset, limit int) ([]domain.Executor, error)
	Save(ctx context.Context, exec domain.Executor) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type executorService struct {
	repo repository.ExecutorRepository
}

func NewExecutorService(repo repository.ExecutorRepository) ExecutorService {
	return &executorService{
		repo: repo,
	}
}

func (svc *executorService) List(ctx context.Context, offset, limit int) ([]domain.Executor, error) {
	execs, err := svc.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return execs, nil
}

func (svc *executorService) Save(ctx context.Context, exec domain.Executor) (int64, error) {
	if exec.Id > 0 {
		err := svc.update(ctx, exec)
		return exec.Id, err
	}
	return svc.create(ctx, exec)
}

func (svc *executorService) create(ctx context.Context, exec domain.Executor) (int64, error) {
	return svc.repo.Create(ctx, exec)
}

func (svc *executorService) update(ctx context.Context, exec domain.Executor) error {
	return svc.repo.Update(ctx, exec)
}

func (svc *executorService) Delete(ctx context.Context, id int64) error {
	return svc.repo.Delete(ctx, id)
}

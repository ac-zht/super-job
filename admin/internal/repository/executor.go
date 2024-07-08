package repository

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/repository/dao"
	"strings"
)

type ExecutorRepository interface {
	List(ctx context.Context, offset, limit int) ([]domain.Executor, error)
	Create(ctx context.Context, exec domain.Executor) (int64, error)
	Update(ctx context.Context, exec domain.Executor) error
	Delete(ctx context.Context, id int64) error
}

type executorRepository struct {
	dao dao.ExecutorDAO
}

func NewExecutorRepository(dao dao.ExecutorDAO) ExecutorRepository {
	return &executorRepository{
		dao: dao,
	}
}

func (repo *executorRepository) List(ctx context.Context, offset, limit int) ([]domain.Executor, error) {
	execs, err := repo.dao.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Executor, domain.Executor](execs, func(idx int, src dao.Executor) domain.Executor {
		return repo.toDomain(src)
	}), nil
}

func (repo *executorRepository) Create(ctx context.Context, exec domain.Executor) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(exec))
}

func (repo *executorRepository) Update(ctx context.Context, exec domain.Executor) error {
	return repo.dao.Update(ctx, repo.toEntity(exec))
}

func (repo *executorRepository) Delete(ctx context.Context, id int64) error {
	return repo.dao.Delete(ctx, id)
}

func (repo *executorRepository) toEntity(e domain.Executor) dao.Executor {
	return dao.Executor{
		Id:    e.Id,
		Name:  e.Name,
		Hosts: strings.Join(e.Hosts, ","),
		Ctime: e.Ctime,
		Utime: e.Utime,
	}
}

func (repo *executorRepository) toDomain(e dao.Executor) domain.Executor {
	return domain.Executor{
		Id:    e.Id,
		Name:  e.Name,
		Hosts: strings.Split(e.Hosts, ","),
		Ctime: e.Ctime,
		Utime: e.Utime,
	}
}

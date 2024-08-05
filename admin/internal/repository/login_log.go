package repository

import (
	"context"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
	"time"
)

type LoginLogRepository interface {
	Create(ctx context.Context, l domain.LoginLog) (int64, error)
	List(ctx context.Context, offset, limit int) ([]domain.LoginLog, error)
	Total(ctx context.Context) (int64, error)
}

type loginLogRepository struct {
	dao dao.LoginLogDAO
}

func NewLoginLogRepository(dao dao.LoginLogDAO) LoginLogRepository {
	return &loginLogRepository{dao: dao}
}

func (repo *loginLogRepository) Create(ctx context.Context, l domain.LoginLog) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(l))
}

func (repo *loginLogRepository) List(ctx context.Context, offset, limit int) ([]domain.LoginLog, error) {
	logs, err := repo.dao.List(ctx, offset, limit)
	if err != nil {
		return []domain.LoginLog{}, err
	}
	return slice.Map[dao.LoginLog, domain.LoginLog](logs, func(idx int, src dao.LoginLog) domain.LoginLog {
		return repo.toDomain(src)
	}), nil
}

func (repo *loginLogRepository) Total(ctx context.Context) (int64, error) {
	return repo.dao.Total(ctx)
}

func (repo *loginLogRepository) toEntity(log domain.LoginLog) dao.LoginLog {
	return dao.LoginLog{
		Id:       log.Id,
		Username: log.Username,
		Ip:       log.Ip,
	}
}

func (repo *loginLogRepository) toDomain(log dao.LoginLog) domain.LoginLog {
	return domain.LoginLog{
		Id:       log.Id,
		Username: log.Username,
		Ip:       log.Ip,
		Ctime:    time.UnixMilli(log.Ctime),
	}
}

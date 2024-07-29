package repository

import (
	"context"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ecodeclub/ekit/slice"
)

type UserRepository interface {
	List(ctx context.Context, offset, limit int) ([]domain.User, error)
	GetById(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, j domain.User) (int64, error)
	Update(ctx context.Context, task domain.User) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	dao dao.UserDAO
}

func (repo *userRepository) List(ctx context.Context, offset, limit int) ([]domain.User, error) {
	users, err := repo.dao.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.User, domain.User](users, func(idx int, src dao.User) domain.User {
		return repo.toDomain(src)
	}), nil
}

func (repo *userRepository) GetById(ctx context.Context, id int64) (domain.User, error) {
	user, err := repo.dao.GetById(ctx, id)
	return repo.toDomain(user), err
}

func (repo *userRepository) Create(ctx context.Context, u domain.User) (int64, error) {
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *userRepository) Update(ctx context.Context, u domain.User) error {
	return repo.dao.Update(ctx, repo.toEntity(u))
}

func (repo *userRepository) Delete(ctx context.Context, id int64) error {
	return repo.dao.Delete(ctx, id)
}

func (repo *userRepository) toEntity(u domain.User) dao.User {
	return dao.User{
		Id:       u.Id,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Status:   u.Status,
	}
}

func (repo *userRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		Status:   u.Status,
		IsAdmin:  u.IsAdmin,
	}
}

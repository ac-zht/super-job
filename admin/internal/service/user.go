package service

import (
	"context"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository"
)

type UserService interface {
	List(ctx context.Context, offset, limit int) ([]domain.User, error)
	GetById(ctx context.Context, id int64) (domain.User, error)
	Save(ctx context.Context, u domain.User) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) List(ctx context.Context, offset, limit int) ([]domain.User, error) {
	return svc.repo.List(ctx, offset, limit)
}

func (svc *userService) GetById(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.GetById(ctx, id)
}

func (svc *userService) Save(ctx context.Context, u domain.User) (int64, error) {
	if u.Id > 0 {
		err := svc.update(ctx, u)
		return u.Id, err
	}
	return svc.create(ctx, u)
}

func (svc *userService) create(ctx context.Context, u domain.User) (int64, error) {
	return svc.repo.Create(ctx, u)
}

func (svc *userService) update(ctx context.Context, u domain.User) error {
	return svc.repo.Update(ctx, u)
}

func (svc *userService) Delete(ctx context.Context, id int64) error {
	return svc.repo.Delete(ctx, id)
}

package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository"
	"github.com/gin-gonic/gin"
)

var ErrUserDuplicate = repository.ErrUserDuplicate

type UserService interface {
	List(ctx context.Context, offset, limit int) ([]domain.User, error)
	Count(ctx context.Context) (int64, error)
	GetById(ctx context.Context, id int64) (domain.User, error)
	ValidateLogin(ctx *gin.Context, username, password string) (domain.User, error)
	Save(ctx context.Context, u domain.User) (int64, error)
	Delete(ctx context.Context, id int64) error
}

type userService struct {
	userRepo     repository.UserRepository
	loginLogRepo repository.LoginLogRepository
}

func NewUserService(userRepo repository.UserRepository, loginLogRepo repository.LoginLogRepository) UserService {
	return &userService{
		userRepo:     userRepo,
		loginLogRepo: loginLogRepo,
	}
}

func (svc *userService) List(ctx context.Context, offset, limit int) ([]domain.User, error) {
	return svc.userRepo.List(ctx, offset, limit)
}

func (svc *userService) Count(ctx context.Context) (int64, error) {
	return svc.userRepo.Count(ctx)
}

func (svc *userService) GetById(ctx context.Context, id int64) (domain.User, error) {
	return svc.userRepo.GetById(ctx, id)
}

func (svc *userService) ValidateLogin(ctx *gin.Context, username, password string) (domain.User, error) {
	user, ok := svc.userRepo.MatchByUsernameAndPassword(ctx, username, password)
	if !ok {
		return domain.User{}, errors.New("用户名或密码错误")
	}
	_, err := svc.loginLogRepo.Create(ctx, domain.LoginLog{
		Username: username,
		Ip:       ctx.RemoteIP(),
	})
	if err != nil {
		return domain.User{}, errors.New(fmt.Sprintf("记录用户登录日志失败#%v", err))
	}
	return user, nil
}

func (svc *userService) Save(ctx context.Context, u domain.User) (int64, error) {
	if u.Id > 0 {
		err := svc.update(ctx, u)
		return u.Id, err
	}
	return svc.create(ctx, u)
}

func (svc *userService) create(ctx context.Context, u domain.User) (int64, error) {
	return svc.userRepo.Create(ctx, u)
}

func (svc *userService) update(ctx context.Context, u domain.User) error {
	return svc.userRepo.Update(ctx, u)
}

func (svc *userService) Delete(ctx context.Context, id int64) error {
	return svc.userRepo.Delete(ctx, id)
}

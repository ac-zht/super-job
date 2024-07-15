package service

import (
	"context"
	"github.com/zc-zht/super-job/admin/internal/repository"
)

type InstallService interface {
	Store(ctx context.Context) error
	Status(ctx context.Context) (bool, error)
}

type installService struct {
	settingRepo repository.SettingRepository
}

func (svc *installService) Store(ctx context.Context) error {
	err := svc.settingRepo.InitBasicField(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (svc *installService) Status(ctx context.Context) (bool, error) {
	//TODO implement me
	panic("implement me")
}

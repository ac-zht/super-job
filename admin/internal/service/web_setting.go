package service

import (
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/repository"
)

type WebSettingService interface {
	Read(fn string) (*domain.Setting, error)
	Write(config []string, fn string) error
}

type WebSettingSvc struct {
	repo repository.WebSettingRepository
}

func NewWebSettingService(repo repository.WebSettingRepository) WebSettingService {
	return &WebSettingSvc{
		repo: repo,
	}
}

func (svc *WebSettingSvc) Read(fn string) (*domain.Setting, error) {
	return svc.repo.Read(fn)
}

func (svc *WebSettingSvc) Write(config []string, fn string) error {
	return svc.repo.Write(config, fn)
}

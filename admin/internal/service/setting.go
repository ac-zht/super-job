package service

import (
	"context"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/repository"
)

type SettingService interface {
	Mail(ctx context.Context) (domain.Mail, error)
	Slack(ctx context.Context) (domain.Slack, error)
	WebHook(ctx context.Context) (domain.WebHook, error)
}

type settingService struct {
	repo repository.SettingRepository
}

func (svc *settingService) Mail(ctx context.Context) (domain.Mail, error) {
	return svc.repo.Mail(ctx)
}

func (svc *settingService) Slack(ctx context.Context) (domain.Slack, error) {
	return svc.repo.Slack(ctx)
}

func (svc *settingService) WebHook(ctx context.Context) (domain.WebHook, error) {
	return svc.repo.WebHook(ctx)
}

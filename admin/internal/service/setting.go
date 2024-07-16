package service

import (
	"context"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/repository"
)

type SettingService interface {
	Mail(ctx context.Context) (domain.Mail, error)
	Slack(ctx context.Context) (domain.Slack, error)
	Webhook(ctx context.Context) (domain.Webhook, error)
	UpdateMail(ctx context.Context, mail domain.Mail) (int64, error)
	UpdateSlack(ctx context.Context, setting domain.Slack) (int64, error)
	UpdateWebhook(ctx context.Context, setting domain.Webhook) (int64, error)
	RemoveMailUser(ctx context.Context, id int64) error
	RemoveChannel(ctx context.Context, id int64) error
}

type settingService struct {
	repo repository.SettingRepository
}

func NewSettingService(repo repository.SettingRepository) SettingService {
	return &settingService{
		repo: repo,
	}
}

func (svc *settingService) Mail(ctx context.Context) (domain.Mail, error) {
	return svc.repo.Mail(ctx)
}

func (svc *settingService) Slack(ctx context.Context) (domain.Slack, error) {
	return svc.repo.Slack(ctx)
}

func (svc *settingService) Webhook(ctx context.Context) (domain.Webhook, error) {
	return svc.repo.Webhook(ctx)
}

func (svc *settingService) UpdateMail(ctx context.Context, mail domain.Mail) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *settingService) UpdateSlack(ctx context.Context, setting domain.Slack) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *settingService) UpdateWebhook(ctx context.Context, setting domain.Webhook) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *settingService) RemoveMailUser(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

func (svc *settingService) RemoveChannel(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}

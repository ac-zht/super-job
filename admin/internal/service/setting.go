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
	UpdateMail(ctx context.Context, mailServer domain.MailServer, template string) error
	UpdateSlack(ctx context.Context, slack domain.Slack) error
	UpdateWebhook(ctx context.Context, webhook domain.Webhook) error

	CreateMailUser(ctx context.Context, mailUser domain.MailUser) (int64, error)
	RemoveMailUser(ctx context.Context, id int64) error

	CreateChannel(ctx context.Context, channel domain.Channel) (int64, error)
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

func (svc *settingService) UpdateMail(ctx context.Context, mailServer domain.MailServer, template string) error {
	return svc.repo.UpdateMail(ctx, mailServer, template)
}

func (svc *settingService) UpdateSlack(ctx context.Context, slack domain.Slack) error {
	return svc.repo.UpdateSlack(ctx, slack)
}

func (svc *settingService) UpdateWebhook(ctx context.Context, webhook domain.Webhook) error {
	return svc.repo.UpdateWebhook(ctx, webhook)
}

func (svc *settingService) CreateMailUser(ctx context.Context, mailUser domain.MailUser) (int64, error) {
	return svc.repo.CreateMailUser(ctx, mailUser)
}

func (svc *settingService) CreateChannel(ctx context.Context, channel domain.Channel) (int64, error) {
	return svc.repo.CreateChannel(ctx, channel)
}

func (svc *settingService) RemoveMailUser(ctx context.Context, id int64) error {
	return svc.repo.RemoveMailUser(ctx, id)
}

func (svc *settingService) RemoveChannel(ctx context.Context, id int64) error {
	return svc.repo.RemoveChannel(ctx, id)
}

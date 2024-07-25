package repository

import (
	"context"
	"encoding/json"
	"github.com/ac-zht/super-job/scheduler/internal/domain"
	"github.com/ac-zht/super-job/scheduler/internal/repository/dao"
)

type SettingRepository interface {
	Mail(ctx context.Context) (domain.Mail, error)
	Slack(ctx context.Context) (domain.Slack, error)
	Webhook(ctx context.Context) (domain.Webhook, error)
}

type settingRepository struct {
	dao       dao.SettingDAO
	BasicRows int64
}

func NewSettingRepository(dao dao.SettingDAO) SettingRepository {
	return &settingRepository{
		dao:       dao,
		BasicRows: 6,
	}
}

func (repo *settingRepository) Mail(ctx context.Context) (domain.Mail, error) {
	list, err := repo.dao.FindByKey(ctx, domain.MailCode)
	if err != nil {
		return domain.Mail{}, err
	}
	var (
		mail     domain.Mail
		mailUser domain.MailUser
	)
	for _, v := range list {
		switch v.Key {
		case domain.MailTemplateKey:
			mail.Template = v.Value
		case domain.MailServerKey:
			json.Unmarshal([]byte(v.Value), &mail)
		case domain.MailUserKey:
			json.Unmarshal([]byte(v.Value), &mailUser)
			mailUser.Id = v.Id
			mail.MailUsers = append(mail.MailUsers, mailUser)
		}
	}
	return mail, nil
}

func (repo *settingRepository) Slack(ctx context.Context) (domain.Slack, error) {
	list, err := repo.dao.FindByKey(ctx, domain.SlackCode)
	if err != nil {
		return domain.Slack{}, err
	}
	var (
		slack   domain.Slack
		channel domain.Channel
	)
	for _, v := range list {
		switch v.Key {
		case domain.SlackTemplateKey:
			slack.Template = v.Value
		case domain.SlackUrlKey:
			slack.Url = v.Value
		case domain.SlackChannelKey:
			channel.Id = v.Id
			slack.Channels = append(slack.Channels, channel)
		}
	}
	return slack, nil
}

func (repo *settingRepository) Webhook(ctx context.Context) (domain.Webhook, error) {
	list, err := repo.dao.FindByKey(ctx, domain.WebhookCode)
	if err != nil {
		return domain.Webhook{}, err
	}
	var webhook domain.Webhook
	for _, v := range list {
		switch v.Key {
		case domain.WebTemplateKey:
			webhook.Template = v.Value
		case domain.WebUrlKey:
			webhook.Url = v.Value
		}
	}
	return webhook, nil
}

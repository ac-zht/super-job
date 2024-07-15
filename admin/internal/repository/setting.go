package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/repository/dao"
)

type SettingRepository interface {
	InitBasicField(ctx context.Context) error
	Mail(ctx context.Context) (domain.Mail, error)
	Slack(ctx context.Context) (domain.Slack, error)
	WebHook(ctx context.Context) (domain.WebHook, error)
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

func (repo *settingRepository) InitBasicField(ctx context.Context) error {
	var setting dao.Setting

	basicFields := []struct {
		Code  string
		Key   string
		Value string
	}{
		{
			Code:  domain.MailCode,
			Key:   domain.MailTemplateKey,
			Value: domain.EmailTemplate,
		},
		{
			Code:  domain.MailCode,
			Key:   domain.MailServerKey,
			Value: "",
		},
		{
			Code:  domain.SlackCode,
			Key:   domain.SlackTemplateKey,
			Value: domain.SlackTemplate,
		},
		{
			Code:  domain.SlackCode,
			Key:   domain.SlackUrlKey,
			Value: "",
		},
		{
			Code:  domain.WebhookCode,
			Key:   domain.WebTemplateKey,
			Value: domain.WebhookTemplate,
		},
		{
			Code:  domain.WebhookCode,
			Key:   domain.WebUrlKey,
			Value: "",
		},
	}
	var (
		id  int64
		err error
	)
	for _, v := range basicFields {
		setting.Code = v.Code
		setting.Key = v.Key
		setting.Value = v.Value
		id, err = repo.dao.Insert(ctx, setting)
		if err != nil {
			return err
		}
	}
	if id != repo.BasicRows {
		return errors.New("init rows not meeting expectations")
	}
	return nil
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

func (repo *settingRepository) WebHook(ctx context.Context) (domain.WebHook, error) {
	list, err := repo.dao.FindByKey(ctx, domain.WebhookCode)
	if err != nil {
		return domain.WebHook{}, err
	}
	var webhook domain.WebHook
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

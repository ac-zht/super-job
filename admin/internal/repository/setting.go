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
	Webhook(ctx context.Context) (domain.Webhook, error)

	UpdateMail(ctx context.Context, mailServer domain.MailServer, template string) error
	UpdateSlack(ctx context.Context, slack domain.Slack) error
	UpdateWebhook(ctx context.Context, webhook domain.Webhook) error

	CreateMailUser(ctx context.Context, mailUser domain.MailUser) (int64, error)
	RemoveMailUser(ctx context.Context, id int64) error

	CreateChannel(ctx context.Context, channel domain.Channel) (int64, error)
	RemoveChannel(ctx context.Context, id int64) error
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

func (repo *settingRepository) UpdateMail(ctx context.Context, mailServer domain.MailServer, template string) error {
	serverJsonByte, _ := json.Marshal(mailServer)
	serverJson := string(serverJsonByte)
	err := repo.dao.UpdateByCodeKey(ctx, domain.MailCode, domain.MailServerKey, serverJson)
	if err != nil {
		return err
	}
	err = repo.dao.UpdateByCodeKey(ctx, domain.MailCode, domain.MailTemplateKey, template)
	if err != nil {
		return err
	}
	return nil
}

func (repo *settingRepository) UpdateSlack(ctx context.Context, slack domain.Slack) error {
	err := repo.dao.UpdateByCodeKey(ctx, domain.SlackCode, domain.SlackUrlKey, slack.Url)
	if err != nil {
		return err
	}
	err = repo.dao.UpdateByCodeKey(ctx, domain.SlackCode, domain.SlackTemplateKey, slack.Template)
	if err != nil {
		return err
	}
	return nil
}

func (repo *settingRepository) UpdateWebhook(ctx context.Context, webhook domain.Webhook) error {
	err := repo.dao.UpdateByCodeKey(ctx, domain.WebhookCode, domain.WebUrlKey, webhook.Url)
	if err != nil {
		return err
	}
	err = repo.dao.UpdateByCodeKey(ctx, domain.WebhookCode, domain.WebTemplateKey, webhook.Template)
	if err != nil {
		return err
	}
	return nil
}

func (repo *settingRepository) CreateMailUser(ctx context.Context, mailUser domain.MailUser) (int64, error) {
	mailUserJsonByte, _ := json.Marshal(mailUser)
	return repo.dao.Insert(ctx, dao.Setting{
		Code:  domain.MailCode,
		Key:   domain.MailUserKey,
		Value: string(mailUserJsonByte),
	})
}

func (repo *settingRepository) RemoveMailUser(ctx context.Context, id int64) error {
	return repo.dao.Delete(ctx, id)
}

func (repo *settingRepository) CreateChannel(ctx context.Context, channel domain.Channel) (int64, error) {
	return repo.dao.Insert(ctx, dao.Setting{
		Code:  domain.SlackCode,
		Key:   domain.SlackChannelKey,
		Value: channel.Name,
	})
}

func (repo *settingRepository) RemoveChannel(ctx context.Context, id int64) error {
	return repo.dao.Delete(ctx, id)
}

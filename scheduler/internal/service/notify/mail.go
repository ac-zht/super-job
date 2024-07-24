package notify

import (
	"context"
	"fmt"
	"github.com/ac-zht/super-job/scheduler/internal/domain"
	"github.com/ac-zht/super-job/scheduler/internal/repository"
	"github.com/ac-zht/super-job/scheduler/pkg/utils"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
	"strconv"
	"strings"
	"time"
)

type Mail struct {
	repo          repository.SettingRepository
	retryTimes    uint8
	retryInterval time.Duration
}

func NewMail(repo repository.SettingRepository, retryTimes uint8, retryInterval time.Duration) *Mail {
	return &Mail{
		repo:          repo,
		retryTimes:    retryTimes,
		retryInterval: retryInterval,
	}
}

func (m *Mail) Send(ctx context.Context, msg Message) {
	mail, err := m.repo.Mail(ctx)
	if err != nil {
		zap.L().Error("#mail#从数据库获取mail配置失败", zap.Error(err))
		return
	}
	if mail.Host == "" {
		zap.L().Error("#mail#Host为空")
		return
	}
	if mail.Port == 0 {
		zap.L().Error("#mail#Port为空")
		return
	}
	if mail.User == "" {
		zap.L().Error("#mail#User为空")
		return
	}
	if mail.Password == "" {
		zap.L().Error("#mail#Password为空")
		return
	}
	msg["content"] = parseNotifyTemplate(mail.Template, msg)
	toUsers := m.getActiveMailUsers(mail, msg)
	m.send(mail, toUsers, msg)
}

func (m *Mail) send(mail domain.Mail, toUsers []string, msg Message) {
	body := msg["content"].(string)
	body = strings.Replace(body, "\n", "<br>", -1)
	gomailMessage := gomail.NewMessage()
	gomailMessage.SetHeader("From", mail.User)
	gomailMessage.SetHeader("To", toUsers...)
	gomailMessage.SetHeader("Subject", "super job-定时任务通知")
	gomailMessage.SetBody("text/html", body)
	mailer := gomail.NewDialer(mail.Host, mail.Port, mail.User, mail.Password)
	var i uint8
	for {
		err := mailer.DialAndSend(gomailMessage)
		if err == nil {
			break
		}
		i++
		time.Sleep(m.retryInterval)
		if i < m.retryTimes {
			zap.L().Error(fmt.Sprintf("mail#发送消息失败#%s#消息内容-%s", err.Error(), msg["content"]))
			continue
		}
		break
	}
}

func (m *Mail) getActiveMailUsers(mail domain.Mail, msg Message) []string {
	taskReceiverIds := strings.Split(msg["task_receiver_id"].(string), ",")
	users := []string{}
	for _, v := range mail.MailUsers {
		if utils.InStringSlice(taskReceiverIds, strconv.Itoa(int(v.Id))) {
			users = append(users, v.Email)
		}
	}
	return users
}

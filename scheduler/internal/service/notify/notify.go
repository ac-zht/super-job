package notify

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ac-zht/super-job/scheduler/internal/domain"
	"go.uber.org/zap"
	"html/template"
	"time"
)

type Message map[string]interface{}

type Notifiable interface {
	Send(ctx context.Context, msg Message)
}

type Service struct {
	mail    Mail
	slack   Slack
	webhook Webhook
	queue   chan Message
}

func NewService(capacity int, mail Mail, slack Slack, webhook Webhook) *Service {
	return &Service{
		queue:   make(chan Message, capacity),
		mail:    mail,
		slack:   slack,
		webhook: webhook,
	}
}

func (s *Service) Push(msg Message) {
	s.queue <- msg
}

func (s *Service) Run(ctx context.Context) {
	for msg := range s.queue {
		taskType, taskTypeOk := msg["task_type"]
		_, taskReceiverIdOk := msg["task_receiver_id"]
		_, nameOk := msg["name"]
		_, outputOk := msg["output"]
		_, statusOk := msg["status"]
		if !taskTypeOk || !taskReceiverIdOk || !nameOk || !outputOk || !statusOk {
			zap.L().Error(fmt.Sprintf("#notify#参数不完整#%+v", msg))
			continue
		}
		msg["content"] = fmt.Sprintf("============\n============\n============\n任务名称: %s\n状态: %s\n输出:\n %s\n", msg["name"], msg["status"], msg["output"])
		zap.L().Debug(fmt.Sprintf("%+v", msg))
		switch taskType.(domain.NotifyType) {
		case domain.EmailNotification:
			go s.mail.Send(ctx, msg)
		case domain.SlackNotification:
			go s.slack.Send(ctx, msg)
		case domain.WebhookNotification:
			go s.webhook.Send(ctx, msg)
		}
		time.Sleep(time.Second)
	}
}

func parseNotifyTemplate(notifyTemplate string, msg Message) string {
	tmpl, err := template.New("notify").Parse(notifyTemplate)
	if err != nil {
		return fmt.Sprintf("解析通知模板失败: %s", err)
	}
	var buf bytes.Buffer
	_ = tmpl.Execute(&buf, map[string]interface{}{
		"TaskId":   msg["task_id"],
		"TaskName": msg["name"],
		"Status":   msg["status"],
		"Result":   msg["output"],
		"Remark":   msg["remark"],
	})
	return buf.String()
}

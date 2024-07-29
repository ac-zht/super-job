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
	Name() string
	Send(ctx context.Context, msg Message)
}

type Service interface {
	Push(msg Message)
	Run(ctx context.Context)
}

type NtfService struct {
	Channels map[string]Notifiable
	Queue    chan Message
}

func NewService(queueCap int, nts ...Notifiable) Service {
	service := &NtfService{
		Channels: make(map[string]Notifiable),
		Queue:    make(chan Message, queueCap),
	}
	channels := make(map[string]Notifiable)
	for _, v := range nts {
		channels[v.Name()] = v
	}
	service.Channels = channels
	return service
}

func (s *NtfService) Push(msg Message) {
	s.Queue <- msg
}

func (s *NtfService) Run(ctx context.Context) {
	for msg := range s.Queue {
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
		if _, existChannel := s.Channels[taskType.(domain.NotifyType).ToString()]; !existChannel {
			zap.L().Error(fmt.Sprintf("#notify#通知渠道不存在#%+v", msg))
			continue
		}
		go s.Channels[taskType.(domain.NotifyType).ToString()].Send(ctx, msg)
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

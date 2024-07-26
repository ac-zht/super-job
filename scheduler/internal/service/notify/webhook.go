package notify

import (
	"context"
	"fmt"
	"github.com/ac-zht/gotools/option"
	"github.com/ac-zht/super-job/scheduler/internal/repository"
	"time"
)

type Webhook struct {
	repo          repository.SettingRepository
	retryTimes    uint8
	retryInterval time.Duration
}

func NewWebhookNotify(repo repository.SettingRepository,
	opts ...option.Option[Webhook]) Notifiable {
	webhook := &Webhook{
		repo:          repo,
		retryTimes:    3,
		retryInterval: time.Second,
	}
	option.Apply[Webhook](webhook, opts...)
	return webhook
}

func WithWebhookRetryTimes(rt uint8) option.Option[Webhook] {
	return func(service *Webhook) {
		service.retryTimes = rt
	}
}

func WithWebhookRetryInterval(ri time.Duration) option.Option[Webhook] {
	return func(service *Webhook) {
		service.retryInterval = ri
	}
}

func (w *Webhook) Name() string {
	return "webhook"
}

func (w *Webhook) Send(ctx context.Context, msg Message) {
	fmt.Println("webhook notify")
}

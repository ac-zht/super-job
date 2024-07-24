package notify

import (
	"context"
	"fmt"
	"github.com/ac-zht/super-job/scheduler/internal/repository"
	"time"
)

type Webhook struct {
	repo          repository.SettingRepository
	retryTimes    uint8
	retryInterval time.Duration
}

func (w *Webhook) Send(ctx context.Context, msg Message) {
	fmt.Println("webhook notify")
}

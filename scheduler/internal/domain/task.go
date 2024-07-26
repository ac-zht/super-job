package domain

import (
	"context"
	"github.com/robfig/cron/v3"
	"time"
)

type Task struct {
	Id         int64
	ExecId     int64
	Name       string
	Cfg        string
	Expression string
	//任务下一次的执行时间
	NextTime   time.Time
	Status     uint8
	Multi      uint8
	Protocol   TaskProtocol
	HttpMethod HttpMethod

	ExecutorHandler       string
	Command               string
	ExecutorRouteStrategy string

	Timeout       time.Duration
	RetryTimes    int64
	RetryInterval int64

	NotifyStatus     NotifyStatus
	NotifyType       NotifyType
	NotifyReceiverId string
	NotifyKeyword    string

	Creator int64
	Updater int64

	Ctime int64
	Utime int64

	Executor   Executor
	CancelFunc func()
}

func (tk Task) Next(t time.Time) time.Time {
	expr := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom |
		cron.Month | cron.Dow |
		cron.Descriptor)
	s, _ := expr.Parse(tk.Expression)
	return s.Next(t)
}

type TaskProtocol uint8

const (
	TaskHTTP  TaskProtocol = iota + 1 // HTTP
	TaskRPC                           // RPC
	TaskShell                         // 命令
)

func (t TaskProtocol) ToUint8() uint8 {
	return uint8(t)
}

func (t TaskProtocol) ToString() string {
	switch t {
	case TaskHTTP:
		return "HTTP"
	case TaskRPC:
		return "RPC"
	case TaskShell:
		return "SHELL"
	}
	return ""
}

type HttpMethod uint8

const (
	HttpGet HttpMethod = iota + 1
	HttpPost
)

func (h HttpMethod) ToUint8() uint8 {
	return uint8(h)
}

func (h HttpMethod) ToString() string {
	switch h {
	case HttpGet:
		return "GET"
	case HttpPost:
		return "POST"
	}
	return ""
}

type NotifyStatus uint8

const (
	NoNotification NotifyStatus = iota
	FailNotification
	OverNotification
	OverKeywordNotification
)

func (n NotifyStatus) ToUint8() uint8 {
	return uint8(n)
}

type NotifyType uint8

const (
	EmailNotification NotifyType = iota + 1
	SlackNotification
	WebhookNotification
)

func (n NotifyType) ToUint8() uint8 {
	return uint8(n)
}

func (n NotifyType) ToString() string {
	switch n {
	case EmailNotification:
		return "email"
	case SlackNotification:
		return "slack"
	case WebhookNotification:
		return "webhook"
	}
	return ""
}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int64
}

type Handler interface {
	Run(ctx context.Context, task Task, jobUniqueId int64) (string, error)
}

const HttpExecTimeout = 300 * time.Second

const (
	SingleInstanceRun uint8 = iota
	MultiInstanceRun
)

const (
	TaskStatusForbidden uint8 = iota
	TaskStatusWaiting
	TaskStatusRunning
)

const (
	ExecutorFirstRouteStrategy    = "FIRST"
	ExecutorLastRouteStrategy     = "LAST"
	ExecutorPollRouteStrategy     = "POLL"
	ExecutorRandomRouteStrategy   = "RANDOM"
	ExecutorFailoverRouteStrategy = "FAILOVER"
)

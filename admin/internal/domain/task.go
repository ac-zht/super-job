package domain

import (
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

	Timeout       int64
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

func (j Task) Next(t time.Time) time.Time {
	expr := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom |
		cron.Month | cron.Dow |
		cron.Descriptor)
	s, _ := expr.Parse(j.Expression)
	return s.Next(t)
}

type TaskProtocol uint8

const (
	TaskHTTP  TaskProtocol = iota + 1 // HTTP
	TaskRPC                           // RPC
	TaskShell                         // 系统命令
)

func (t TaskProtocol) ToUint8() uint8 {
	return uint8(t)
}

type HttpMethod uint8

const (
	HttpGet HttpMethod = iota + 1
	HttpPost
)

func (t HttpMethod) ToUint8() uint8 {
	return uint8(t)
}

type NotifyStatus uint8

const (
	NoNotification NotifyStatus = iota
	FailNotification
	OverNotification
	OverKeywordNotification
)

func (t NotifyStatus) ToUint8() uint8 {
	return uint8(t)
}

type NotifyType uint8

const (
	EmailNotification NotifyType = iota + 1
	SlackNotification
	WebhookNotification
)

func (t NotifyType) ToUint8() uint8 {
	return uint8(t)
}

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

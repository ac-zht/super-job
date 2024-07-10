package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id         int64
	ExecId     int64
	Name       string
	Cfg        string
	Expression string
	//任务下一次的执行时间
	NextTime   time.Time
	Status     uint8
	Multi      uint8
	Protocol   JobProtocol
	HttpMethod HttpMethod

	Timeout       int64
	RetryTimes    int64
	RetryInterval int64

	Creator int64
	Updater int64

	Ctime int64
	Utime int64

	Executor   Executor
	CancelFunc func()
}

func (j Job) Next(t time.Time) time.Time {
	expr := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom |
		cron.Month | cron.Dow |
		cron.Descriptor)
	s, _ := expr.Parse(j.Expression)
	return s.Next(t)
}

type JobProtocol uint8
type HttpMethod uint8

const (
	TaskHTTP  JobProtocol = iota + 1 // HTTP
	TaskRPC                          // RPC
	TaskShell                        // 系统命令
)

const (
	HttpGet HttpMethod = iota + 1
	HttpPost
)

func (t JobProtocol) ToUint8() uint8 {
	return uint8(t)
}

func (t HttpMethod) ToUint8() uint8 {
	return uint8(t)
}

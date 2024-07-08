package domain

import (
	"github.com/robfig/cron/v3"
	"time"
)

type Job struct {
	Id         int64
	ExecId     int64
	Name       string
	Protocol   JobProtocol
	Cfg        string
	Expression string
	//任务下一次的执行时间
	NextTime time.Time

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

const (
	TaskHTTP  JobProtocol = iota + 1 // HTTP
	TaskRPC                          // RPC
	TaskShell                        // 系统命令
)

func (t JobProtocol) ToUint8() uint8 {
	return uint8(t)
}

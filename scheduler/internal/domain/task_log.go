package domain

import (
	"github.com/ac-zht/super-job/scheduler/internal/repository/dao"
	"time"
)

type TaskLog struct {
	Id     int64
	TaskId int64
	ExecId int64
	Name   string

	Spec          string
	SchedulerAddr string

	Protocol    string
	Command     string
	ExecutorMsg string

	Timeout    int64
	RetryTimes int64

	StartTime time.Time
	EndTime   time.Time
	Status    uint8
	TotalTime int
	Result    string

	NotifyStatus uint8
}

const (
	JobStatusFailure uint8 = iota
	JobStatusRunning
	JobStatusFinish
	JobStatusCancel
)

type CommonMap dao.CommonMap

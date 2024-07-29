package dao

import "context"

type TaskLogDAO interface {
	List(ctx context.Context, taskId int64) ([]TaskLog, error)
	Delete(ctx context.Context, id int64) error
}

type TaskLog struct {
	Id     int64
	TaskId int64  //任务id
	ExecId int64  //执行器id
	Name   string //任务名称

	Spec          string //调度规则方式
	SchedulerAddr string //调度器地址

	Protocol    string //请求协议
	Command     string //shell命令或URL地址
	ExecutorMsg string //执行器信息：名称，执行器注册的地址，本次执行地址，路由策略，执行器任务handler，

	Timeout    int64 //任务执行超时时间
	RetryTimes int64 //失败重试次数

	StartTime int64  //任务开始时间
	EndTime   int64  //任务结束时间
	Status    uint8  //执行状态
	TotalTime int    //执行总时长
	Result    string //执行结果

	NotifyStatus uint8 //通知状态 0:默认 1:无需通知 2:通知成功 3:通知失败
}

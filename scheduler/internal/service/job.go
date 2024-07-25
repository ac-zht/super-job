package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ac-zht/super-job/scheduler/internal/domain"
	"github.com/ac-zht/super-job/scheduler/internal/repository"
	hc "github.com/ac-zht/super-job/scheduler/internal/service/http/client"
	"github.com/ac-zht/super-job/scheduler/internal/service/notify"
	rc "github.com/ac-zht/super-job/scheduler/internal/service/rpc/client"
	"github.com/ac-zht/super-job/scheduler/internal/service/rpc/executor/route"
	pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"
	"github.com/ac-zht/super-job/scheduler/pkg/utils"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

//go:generate mockgen -source=./cron_job.go -package=svcmocks -destination=mocks/cron_job.mocks.go JobService
type JobService interface {
	Preempt(ctx context.Context) (domain.Task, error)
	ResetNextTime(ctx context.Context, task domain.Task) error
	CreateJob(task domain.Task) func(ctx context.Context) error
}

type cronJobService struct {
	taskRepo        repository.TaskRepository
	taskLogRepo     repository.TaskLogRepository
	notifySvc       notify.Service
	refreshInterval time.Duration
}

func NewJobService(taskRepo repository.TaskRepository,
	taskLogRepo repository.TaskLogRepository, notifySvc notify.Service, refreshInterval time.Duration) JobService {
	return &cronJobService{
		taskRepo:        taskRepo,
		taskLogRepo:     taskLogRepo,
		notifySvc:       notifySvc,
		refreshInterval: refreshInterval,
	}
}

func (c *cronJobService) Preempt(ctx context.Context) (domain.Task, error) {
	j, err := c.taskRepo.Preempt(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	ticker := time.NewTicker(c.refreshInterval)
	//续约任务
	go func() {
		for range ticker.C {
			c.refresh(j.Id)
		}
	}()
	j.CancelFunc = func() {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.taskRepo.Release(ctx, j.Id)
		if err != nil {
			zap.L().Error("释放任务失败", zap.Int64("id", j.Id), zap.Error(err))
		}
	}
	return j, nil
}

func (c *cronJobService) CreateJob(task domain.Task) func(ctx context.Context) error {
	handler := c.createHandler(task)
	if handler == nil {
		return nil
	}
	return func(ctx context.Context) error {
		taskLogId := c.beforeExecJob(ctx, task)
		if taskLogId <= 0 {
			return errors.New("任务执行准备失败")
		}
		var exec string
		if task.Protocol == domain.TaskRPC {
			exec = task.ExecutorHandler
		} else {
			exec = task.Command
		}
		zap.L().Info(fmt.Sprintf(" 开始执行任务#%s#命令-%s", task.Name, exec))
		taskResult := c.execJob(ctx, handler, task, taskLogId)
		zap.L().Info(fmt.Sprintf(" 任务完成#%s#命令-%s", task.Name, exec))
		c.afterExecJob(ctx, taskLogId, task, taskResult)
		return nil
	}
}

func (c *cronJobService) ResetNextTime(ctx context.Context, task domain.Task) error {
	t := task.Next(time.Now())
	if !t.IsZero() {
		return c.taskRepo.UpdateNextTime(ctx, task.Id, t)
	}
	return nil
}

func (c *cronJobService) createHandler(task domain.Task) domain.Handler {
	var handler domain.Handler = nil
	switch task.Protocol {
	case domain.TaskHTTP:
		handler = new(HTTPHandler)
	case domain.TaskRPC:
		handler = new(RPCHandler)
	}
	return handler
}

func (c *cronJobService) execJob(ctx context.Context, handler domain.Handler, task domain.Task, jobUniqueId int64) domain.TaskResult {
	execTimes := task.RetryTimes + 1
	var (
		output string
		err    error
	)
	for i := int64(0); i < execTimes; i++ {
		output, err = handler.Run(ctx, task, jobUniqueId)
		if err == nil {
			return domain.TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		if i < execTimes {
			zap.L().Warn(fmt.Sprintf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", task.Id, i, output, err.Error()))
			if task.RetryInterval > 0 {
				time.Sleep(time.Duration(task.RetryInterval) * time.Second)
			} else {
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}
	return domain.TaskResult{Result: output, Err: err, RetryTimes: task.RetryTimes}
}

type HTTPHandler struct{}

func (h *HTTPHandler) Run(ctx context.Context, task domain.Task, jobUniqueId int64) (string, error) {
	if task.Timeout <= 0 || task.Timeout > domain.HttpExecTimeout {
		task.Timeout = domain.HttpExecTimeout
	}
	var resp hc.ResponseWrapper
	httpClient := &hc.HttpClient{
		Url:     task.Command,
		Timeout: task.Timeout,
		Client:  &http.Client{},
	}
	if task.HttpMethod == domain.HttpGet {
		resp = httpClient.Get()
	} else {
		urlFields := strings.Split(task.Command, "?")
		httpClient.Url = urlFields[0]
		var params string
		if len(urlFields) >= 2 {
			params = urlFields[1]
		}
		resp = httpClient.PostParams(params)
	}
	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
	}
	return resp.Body, nil
}

type RPCHandler struct{}

func (h *RPCHandler) Run(ctx context.Context, task domain.Task, jobUniqueId int64) (string, error) {
	taskRequest := new(pb.TaskRequest)
	taskRequest.Type = int32(task.Protocol)
	taskRequest.Id = jobUniqueId
	taskRequest.Timeout = int32(task.Timeout)
	if task.Protocol == domain.TaskRPC {
		taskRequest.Handler = task.ExecutorHandler
	} else {
		taskRequest.Command = task.Command
	}
	var (
		output string
		err    error
	)
	rpcClient := &rc.RpcClient{}
	switch task.ExecutorRouteStrategy {
	case domain.ExecutorFirstRouteStrategy:
		output, err = (&route.FirstRouteReqStrategy{Client: rpcClient}).Call(task.Executor.Hosts, taskRequest)
	case domain.ExecutorLastRouteStrategy:
		output, err = (&route.LastRouteReqStrategy{Client: rpcClient}).Call(task.Executor.Hosts, taskRequest)
	}
	return output, err
}

func (c *cronJobService) beforeExecJob(ctx context.Context, task domain.Task) int64 {
	//检测该任务是否正在被调度
	taskLogId, err := c.createTaskLog(ctx, task, domain.JobStatusRunning)
	if err != nil {
		zap.L().Error(fmt.Sprintf("任务开始执行#写入任务日志失败"), zap.Error(err))
		return 0
	}
	return taskLogId
}

func (c *cronJobService) afterExecJob(ctx context.Context, taskLogId int64, task domain.Task, taskResult domain.TaskResult) {
	defer func() {
		//释放任务
		task.CancelFunc()
	}()
	//设置next_time
	err := c.ResetNextTime(ctx, task)
	if err != nil {
		zap.L().Error(fmt.Sprintf("任务结束#更新下一次的执行失败"), zap.Error(err))
	}
	//更新执行日志
	err = c.updateTaskLog(ctx, taskLogId, taskResult)
	if err != nil {
		zap.L().Error("任务结束#更新任务日志失败", zap.Error(err))
	}
	//通知
	go c.SendNotification(ctx, task, taskResult)
	//执行依赖任务
}

func (c *cronJobService) createTaskLog(ctx context.Context, task domain.Task, status uint8) (int64, error) {
	ip, _ := utils.GetLocalIP()
	var protocol string
	if task.Protocol == domain.TaskHTTP {
		protocol = fmt.Sprintf("%s-%s", task.Protocol.ToString(), task.HttpMethod.ToString())
	} else {
		protocol = task.Protocol.ToString()
	}
	taskLog := domain.TaskLog{
		TaskId:        task.Id,
		ExecId:        task.ExecId,
		Name:          task.Name,
		Spec:          task.Expression,
		SchedulerAddr: ip,
		Protocol:      protocol,
		Command:       task.Command,
		ExecutorMsg:   "",
		Timeout:       task.Timeout,
		RetryTimes:    task.RetryTimes,
		StartTime:     time.Now(),
		Status:        status,
	}
	return c.taskLogRepo.Create(ctx, taskLog)
}

func (c *cronJobService) updateTaskLog(ctx context.Context, taskLogId int64, taskResult domain.TaskResult) error {
	result := taskResult.Result
	var status uint8
	if taskResult.Err != nil {
		status = domain.JobStatusFailure
	} else {
		status = domain.JobStatusFinish
	}
	return c.taskLogRepo.UpdateById(ctx, taskLogId, domain.CommonMap{
		"retry_times": taskResult.RetryTimes,
		"status":      status,
		"result":      result,
	})
}

func (c *cronJobService) SendNotification(ctx context.Context, task domain.Task, taskResult domain.TaskResult) {
	if task.NotifyStatus == domain.NoNotification {
		return
	}
	if task.NotifyStatus == domain.FailNotification && taskResult.Err == nil {
		return
	}
	if task.NotifyStatus == domain.OverKeywordNotification && !strings.Contains(taskResult.Result, task.NotifyKeyword) {
		return
	}
	if task.NotifyType != domain.WebhookNotification && task.NotifyReceiverId == "" {
		return
	}
	var statStr string
	if taskResult.Err != nil {
		statStr = "失败"
	} else {
		statStr = "成功"
	}
	msg := notify.Message{
		"task_type":        task.NotifyType,
		"task_receiver_id": task.NotifyReceiverId,
		"name":             task.Name,
		"output":           taskResult.Result,
		"status":           statStr,
		"task_id":          task.Id,
		"remark":           task.Cfg,
	}
	c.notifySvc.Push(msg)
}

func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.taskRepo.UpdateUtime(ctx, id)
	if err != nil {
		zap.L().Error("续约失败", zap.Int64("jid", id), zap.Error(err))
	}
}

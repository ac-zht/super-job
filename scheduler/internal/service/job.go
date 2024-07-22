package service

import (
	"context"
	"fmt"
	"github.com/ac-zht/super-job/scheduler/internal/domain"
	"github.com/ac-zht/super-job/scheduler/internal/repository"
	pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"
	"github.com/ac-zht/super-job/scheduler/pkg/logger"
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
	repo            repository.TaskRepository
	l               logger.Logger
	refreshInterval time.Duration
}

func NenJobService(
	repo repository.TaskRepository,
	l logger.Logger) JobService {
	return &cronJobService{
		repo:            repo,
		l:               l,
		refreshInterval: time.Second * 10,
	}
}

func (c *cronJobService) Preempt(ctx context.Context) (domain.Task, error) {
	j, err := c.repo.Preempt(ctx)
	if err != nil {
		return domain.Task{}, err
	}
	ticker := time.NewTicker(c.refreshInterval)
	go func() {
		for range ticker.C {
			c.refresh(j.Id)
		}
	}()
	j.CancelFunc = func() {
		ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := c.repo.Release(ctx, j.Id)
		if err != nil {
			c.l.Error("释放任务失败",
				logger.Error(err),
				logger.Int64("id", j.Id))
		}
	}
	return j, nil
}

func (c *cronJobService) ResetNextTime(ctx context.Context, task domain.Task) error {
	t := task.Next(time.Now())
	if !t.IsZero() {
		return c.repo.UpdateNextTime(ctx, task.Id, t)
	}
	return nil
}

func (c *cronJobService) CreateJob(task domain.Task) func(ctx context.Context) error {
	handler := c.createHandler(task)
	if handler == nil {
		return nil
	}
	return func(ctx context.Context) error {
		c.beforeExecJob(task)
		c.execJob(handler, task, 1)
		c.afterExecJob(task)
		return nil
	}
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

func (c *cronJobService) execJob(handler domain.Handler, task domain.Task, jobUniqueId int64) domain.TaskResult {
	execTimes := task.RetryTimes + 1
	var (
		output string
		err    error
	)
	for i := int8(0); i < execTimes; i++ {
		output, err = handler.Run(task, jobUniqueId)
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

func (c *cronJobService) beforeExecJob(task domain.Task) {
	//创建执行日志
}

func (c *cronJobService) afterExecJob(task domain.Task) {
	//设置next_time
	//更新执行日志
	//通知
	//执行依赖任务
}

type HTTPHandler struct{}

func (h *HTTPHandler) Run(task domain.Task, jobUniqueId int64) (string, error) {
	if task.Timeout <= 0 || task.Timeout > domain.HttpExecTimeout {
		task.Timeout = domain.HttpExecTimeout
	}
	var resp ResponseWrapper
	httpClient := &HttpClient{
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

func (h *RPCHandler) Run(task domain.Task, jobUniqueId int64) (string, error) {
	taskRequest := new(pb.TaskRequest)
	taskRequest.Timeout = int32(task.Timeout)
	taskRequest.Command = task.Command
	taskRequest.Id = jobUniqueId
	switch task.ExecutorRouteStrategy {
	case domain.ExecutorFirstRouteStrategy:

	case domain.ExecutorLastRouteStrategy:
	}

}

func (c *cronJobService) refresh(id int64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := c.repo.UpdateUtime(ctx, id)
	if err != nil {
		c.l.Error("续约失败",
			logger.Int64("jid", id),
			logger.Error(err))
	}
}

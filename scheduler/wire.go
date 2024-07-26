package main

import (
	"github.com/ac-zht/gotools/option"
	"github.com/ac-zht/gotools/pool"
	"github.com/ac-zht/super-job/scheduler/internal/repository"
	"github.com/ac-zht/super-job/scheduler/internal/repository/dao"
	"github.com/ac-zht/super-job/scheduler/internal/service"
	"github.com/ac-zht/super-job/scheduler/internal/service/http/client"
	"github.com/ac-zht/super-job/scheduler/internal/service/notify"
	"github.com/ac-zht/super-job/scheduler/internal/web"
	"github.com/ac-zht/super-job/scheduler/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"net/http"
)

func InjectHttpClient() *client.HttpClient {
	return &client.HttpClient{
		Url:     "",
		Timeout: 5,
		Client:  &http.Client{},
		Req:     &http.Request{},
	}
}

//func NewNotifyService() notify.Service {
//    wire.NewSet(
//        dao.NewSettingDAO,
//        repository.NewSettingRepository,
//
//        notify.NewMailNotify,
//        InjectHttpClient,
//        notify.NewSlackNotify,
//        notify.NewWebhookNotify,
//    )
//    wire.Struct(&notify.NtfService{Queue: make(chan notify.Message, 100)}, "Channels")
//    return &notify.NtfService{}
//}

var notifyServiceProvider = wire.NewSet(
	dao.NewSettingDAO,
	repository.NewSettingRepository,
	notify.NewMailNotify,
	InjectHttpClient,
	notify.NewSlackNotify,
	notify.NewWebhookNotify,
)

func NewOnDemandBlockTaskPool() *pool.OnDemandBlockTaskPool {
	quickPool, _ := pool.NewOnDemandBlockTaskPool(5, 10)
	return quickPool
}

func NewCronJobServiceOption() []option.Option[service.CronJobService] {
	return []option.Option[service.CronJobService]{func(s *service.CronJobService) {}}
}

func NewSchedulerOption() []option.Option[web.Scheduler] {
	return []option.Option[web.Scheduler]{func(s *web.Scheduler) {}}
}

func InitScheduler() *gin.Engine {
	wire.Build(
		ioc.InitDB,

		dao.NewTaskDAO,
		dao.NewTaskLogDAO,
		repository.NewTaskRepository,
		repository.NewTaskLogRepository,

		//NewNotifyService,
		notifyServiceProvider,
		NewCronJobServiceOption,
		service.NewJobService,

		NewOnDemandBlockTaskPool,
		NewSchedulerOption,
		web.NewScheduler,

		ioc.InitWebServer,
	)
	wire.Struct(&notify.NtfService{Queue: make(chan notify.Message, 100)}, "Channels")
	return new(gin.Engine)
}

package handler

import taskSvc "github.com/ac-zht/super-job/executor/proto"

type JobHandler interface {
	Run() (*taskSvc.TaskResponse, error)
}

var Handlers = make(map[string]JobHandler)

func RegisterJobHandler(name string, h JobHandler) {
	Handlers[name] = h
}

func init() {
	demo := DemoJobHandler{}
	RegisterJobHandler("demoJobHandler", demo)
}

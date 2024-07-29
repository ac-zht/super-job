package handler

import taskSvc "github.com/ac-zht/super-job/executor/proto"

type JobHandler interface {
	Name() string
	Run() (*taskSvc.TaskResponse, error)
}

var Handlers = make(map[string]JobHandler)

func RegisterJobHandler(h JobHandler) {
	Handlers[h.Name()] = h
}

func init() {
	RegisterJobHandler(
		&DemoJobHandler{},
	)
}

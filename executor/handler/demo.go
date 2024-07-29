package handler

import (
	"fmt"
	taskSvc "github.com/ac-zht/super-job/executor/proto"
)

type DemoJobHandler struct {
}

func (d *DemoJobHandler) Name() string {
	return "demoJobHandler"
}

func (d *DemoJobHandler) Run() (*taskSvc.TaskResponse, error) {
	fmt.Println("DemoJobHandler is executing")
	return &taskSvc.TaskResponse{
		Output: "success",
		Error:  "0",
	}, nil
}

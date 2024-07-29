package main

import (
	"context"
	"fmt"
	"github.com/ac-zht/super-job/executor/handler"
	taskSvc "github.com/ac-zht/super-job/executor/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

const (
	HandlerRequest = iota + 2
	CommandRequest
)

type Server struct {
	taskSvc.UnimplementedTaskServer
}

func (s *Server) Run(ctx context.Context, request *taskSvc.TaskRequest) (*taskSvc.TaskResponse, error) {
	if request.Type == HandlerRequest {
		return handler.Handlers[request.Handler].(handler.JobHandler).Run()
	}
	return &taskSvc.TaskResponse{
		Output: request.Command,
		Error:  "0",
	}, nil
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	defer logger.Sync()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()
	grpcServer := grpc.NewServer()
	taskSvc.RegisterTaskServer(grpcServer, &Server{})
	listen, err := net.Listen("tcp", ":9300")
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("failed to listen: %v", err))
	}
	zap.L().Info("服务启动中...")
	if err = grpcServer.Serve(listen); err != nil {
		zap.L().Fatal(fmt.Sprintf("failed to serve: %v", err))
	}
}

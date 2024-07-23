package client

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

var (
	taskMap        sync.Map
	errUnavailable = errors.New("无法连接远程服务器")
)

type RpcClient struct {
}

func (r *RpcClient) Stop(addr string, id int64) {
	key := generateTaskUniqueKey(addr, id)
	cancel, ok := taskMap.Load(key)
	if !ok {
		return
	}
	cancel.(context.CancelFunc)()
}

func (r *RpcClient) Exec(addr string, taskReq *pb.TaskRequest) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			zap.L().Error(fmt.Sprintf("panic#rpc/client.go:Exec#$%s", err.(error).Error()))
		}
	}()
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewTaskClient(conn)
	//超时限制一天
	if taskReq.Timeout <= 0 || taskReq.Timeout > 86400 {
		taskReq.Timeout = 86400
	}

	timeout := time.Duration(taskReq.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	taskUniqueKey := generateTaskUniqueKey(addr, taskReq.Id)
	taskMap.Store(taskUniqueKey, cancel)
	defer taskMap.Delete(taskUniqueKey)

	resp, err := c.Run(ctx, taskReq)
	if err != nil {
		return parseGRPCError(err)
	}
	if resp.Error == "" {
		return resp.Output, nil
	}
	return resp.Output, errors.New(resp.Error)
}

func generateTaskUniqueKey(addr string, id int64) string {
	return fmt.Sprintf("%s:%d", addr, id)
}

func parseGRPCError(err error) (string, error) {
	switch status.Code(err) {
	case codes.Unavailable:
		return "", errUnavailable
	case codes.DeadlineExceeded:
		return "", errors.New("执行超时，强制结束")
	case codes.Canceled:
		return "", errors.New("手动停止")
	}
	return "", err
}

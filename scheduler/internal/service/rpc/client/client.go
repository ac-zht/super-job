package client

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type RpcClient struct {
}

func (r *RpcClient) Stop(addr string, id int64) {

}

func (r *RpcClient) Exec(addr string, taskReq *pb.TaskRequest) (string, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewTaskClient(conn)
	timeout := time.Duration(taskReq.Timeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resp, err := c.Run(ctx, taskReq)
	if err != nil {
		return "", err
	}
	if resp.Error == "" {
		return resp.Output, nil
	}
	return resp.Output, errors.New(resp.Error)
}

func generateTaskUniqueKey(addr string, id int64) string {
	return fmt.Sprintf("%s:%d", addr, id)
}

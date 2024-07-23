package route

import (
	"github.com/ac-zht/super-job/scheduler/internal/service/rpc/client"
	pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"
)

type FirstRouteReqStrategy struct {
	Client *client.RpcClient
}

func (f *FirstRouteReqStrategy) Call(hosts []string, taskReq *pb.TaskRequest) (string, error) {
	return f.Client.Exec(hosts[0], taskReq)
}

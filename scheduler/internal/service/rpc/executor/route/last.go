package route

import (
	"github.com/ac-zht/super-job/scheduler/internal/service/rpc/client"
	pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"
)

type LastRouteReqStrategy struct {
	Client *client.RpcClient
}

func (f *LastRouteReqStrategy) Call(hosts []string, taskReq *pb.TaskRequest) (string, error) {
	return f.Client.Exec(hosts[len(hosts)-1], taskReq)
}

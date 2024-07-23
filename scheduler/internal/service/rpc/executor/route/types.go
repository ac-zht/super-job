package route

import pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"

type ReqStrategy interface {
	Call(hosts []string, taskReq *pb.TaskRequest) (string, error)
}

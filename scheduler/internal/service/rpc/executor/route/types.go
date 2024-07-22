package route

import pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"

type Strategy interface {
	Call(addrs string, request *pb.TaskRequest) (string, error)
}

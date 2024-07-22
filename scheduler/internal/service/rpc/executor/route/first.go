package route

import (
	pb "github.com/ac-zht/super-job/scheduler/internal/service/rpc/proto"
	"strings"
)

type FirstRouteStrategy struct {
}

func (f *FirstRouteStrategy) Call(addrs string, request *pb.TaskRequest) (string, error) {
	addrList := strings.Split(addrs, ",")
	//请求
	//返回
}

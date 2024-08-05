package web

type TaskVo struct {
	Id         int64  `json:"id"`
	Executor   string `json:"executor"`
	Name       string `json:"name"`
	Cfg        string `json:"cfg"`
	Expression string `json:"expression"`
	Status     uint8  `json:"status"`
	Multi      uint8  `json:"multi"`

	Protocol        uint8  `json:"protocol"`
	HttpMethod      uint8  `json:"http_method"`
	ExecutorHandler string `json:"executor_handler"`
	Command         string `json:"command"`

	Timeout       int64 `json:"timeout"`
	RetryTimes    int64 `json:"retry_times"`
	RetryInterval int64 `json:"retry_interval"`

	Ctime    int64 `json:"ctime"`
	NextTime int64 `json:"next_time"`
}

type TaskDetail struct {
	Id         int64  `json:"id"`
	ExecId     int64  `json:"exec_id"`
	Name       string `json:"name"`
	Cfg        string `json:"cfg"`
	Expression string `json:"expression"`
	Status     uint8  `json:"status"`
	Multi      uint8  `json:"multi"`

	Protocol        uint8  `json:"protocol"`
	HttpMethod      uint8  `json:"httpMethod"`
	ExecutorHandler string `json:"executor_handler"`
	Command         string `json:"command"`

	Timeout       int64 `json:"timeout"`
	RetryTimes    int64 `json:"retry_times"`
	RetryInterval int64 `json:"retry_interval"`

	NotifyStatus     uint8  `json:"notify_status"`
	NotifyType       uint8  `json:"notify_type"`
	NotifyReceiverId string `json:"notify_receiver_id"`
	NotifyKeyword    string `json:"notify_keyword"`
}

type ExecutorVo struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Hosts string `json:"hosts"`
	Ctime int64  `json:"ctime"`
	Utime int64  `json:"utime"`
}

type ExecutorBrief struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type LoginResp struct {
	Token    string `json:"token"`
	Uid      int64  `json:"uid"`
	Username string `json:"username"`
	IsAdmin  uint8  `json:"is_admin"`
}

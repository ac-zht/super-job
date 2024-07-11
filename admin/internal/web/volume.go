package web

type JobVo struct {
	Id         int64  `json:"id"`
	Executor   string `json:"executor"`
	Name       string `json:"name"`
	Cfg        string `json:"cfg"`
	Expression string `json:"expression"`
	Protocol   uint8  `json:"protocol"`
	Status     uint8  `json:"status"`
	Multi      uint8  `json:"multi"`
	HttpMethod uint8  `json:"httpMethod"`

	Timeout       int64 `json:"timeout"`
	RetryTimes    int64 `json:"retry_times"`
	RetryInterval int64 `json:"retry_interval"`

	Ctime    int64 `json:"ctime"`
	NextTime int64 `json:"next_time"`
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

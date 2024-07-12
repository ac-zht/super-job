package web

type JobEditReq struct {
	Id               int64  `json:"id"`
	Name             string `json:"name"`
	ExecId           int64  `json:"exec_id"`
	Cfg              string `json:"cfg"`
	Expression       string `json:"expression"`
	Status           uint8  `json:"status"`
	Multi            uint8  `json:"multi"`
	Protocol         uint8  `json:"protocol"`
	HttpMethod       uint8  `json:"http_method"`
	ExecutorHandler  string `json:"executor_handler"`
	Command          string `json:"command"`
	Timeout          int64  `json:"timeout"`
	RetryTimes       int64  `json:"retry_times"`
	RetryInterval    int64  `json:"retry_interval"`
	NotifyStatus     uint8  `json:"notify_status"`
	NotifyType       uint8  `json:"notify_type"`
	NotifyReceiverId string `json:"notify_receiver_id"`
	NotifyKeyword    string `json:"notify_key_word"`
}

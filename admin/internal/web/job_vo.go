package web

type JobVo struct {
	Id       int64  `json:"id"`
	Executor string `json:"executor"`
	Name     string `json:"name"`
	Protocol uint8  `json:"protocol"`
	Cfg      string `json:"cfg"`
	NextTime int64  `json:"next_time"`
}

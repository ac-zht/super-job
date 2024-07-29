package web

type TaskEditReq struct {
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

type MailEditReq struct {
	Host     string `json:"host,omitempty" validate:"required"`
	Port     int    `json:"port,omitempty" validate:"required"`
	User     string `json:"user,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
	Template string `json:"template,omitempty" validate:"required"`
}

type UserEditReq struct {
	Id              int64
	Name            string
	Email           string
	Password        string
	ConfirmPassword string `json:"confirm_password"`
	IsAdmin         uint8
	Status          uint8
}

type InstallReq struct {
	DbType               string `json:"db_type,omitempty"`
	DbHost               string `json:"db_host,omitempty"`
	DbPort               int    `json:"db_port,omitempty"`
	DbUsername           string `json:"db_username,omitempty"`
	DbPassword           string `json:"db_password,omitempty"`
	DbName               string `json:"db_name,omitempty"`
	DbTablePrefix        string `json:"db_table_prefix,omitempty"`
	AdminUsername        string `json:"admin_username,omitempty"`
	AdminPassword        string `json:"admin_password,omitempty"`
	ConfirmAdminPassword string `json:"confirm_admin_password,omitempty"`
	AdminEmail           string `json:"admin_email,omitempty"`
}

package domain

type Installation struct {
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

type Setting struct {
	DB struct {
		Host         string
		Port         int
		User         string
		Password     string
		Database     string
		Prefix       string
		Charset      string
		MaxIdleConns int
		MaxOpenConns int
	}
}

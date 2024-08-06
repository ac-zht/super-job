package domain

const DefaultSection = "default"

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
		Engine       string
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
	AllowIps      string
	AppName       string
	ApiKey        string
	ApiSecret     string
	ApiSignEnable bool

	EnableTLS bool
	CAFile    string
	CertFile  string
	KeyFile   string

	AccessTokenKey  string
	RefreshTokenKey string
}

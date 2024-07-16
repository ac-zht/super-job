package domain

type Mail struct {
	Host      string     `json:"host"`
	Port      int        `json:"port"`
	User      string     `json:"user"`
	Password  string     `json:"password"`
	MailUsers []MailUser `json:"mail_users"`
	Template  string     `json:"template"`
}

type MailUser struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Slack struct {
	Url      string    `json:"url"`
	Channels []Channel `json:"channels"`
	Template string    `json:"template"`
}

type Channel struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Webhook struct {
	Url      string `json:"url"`
	Template string `json:"template"`
}

const (
	MailCode        = "mail"
	MailTemplateKey = "template"
	MailServerKey   = "server"
	MailUserKey     = "user"
)

const (
	SlackCode        = "slack"
	SlackTemplateKey = "template"
	SlackUrlKey      = "url"
	SlackChannelKey  = "channel"
)

const (
	WebhookCode    = "webhook"
	WebTemplateKey = "template"
	WebUrlKey      = "url"
)

const EmailTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
备注:    {{.Remark}}
`

const SlackTemplate = `
任务ID:  {{.TaskId}}
任务名称: {{.TaskName}}
状态:    {{.Status}}
执行结果: {{.Result}}
备注:    {{.Remark}}
`

const WebhookTemplate = `
{
  "task_id": "{{.TaskId}}",
  "task_name": "{{.TaskName}}",
  "status": "{{.Status}}",
  "result": "{{.Result}}",
  "remark": "{{.Remark}}",
}
`

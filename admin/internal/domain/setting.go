package domain

type Mail struct {
	Host      string
	Port      int
	User      string
	Password  string
	MailUsers []MailUser
	Template  string
}

type MailUser struct {
	Id       int64
	Username string
	Email    string
}

type Slack struct {
	Url      string
	Channels []Channel
	Template string
}

type Channel struct {
	Id   int64
	Name string
}

type WebHook struct {
	Url      string
	Template string
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

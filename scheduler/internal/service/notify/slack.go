package notify

import (
	"context"
	"fmt"
	"github.com/ac-zht/gotools/option"
	"github.com/ac-zht/super-job/scheduler/internal/domain"
	"github.com/ac-zht/super-job/scheduler/internal/repository"
	"github.com/ac-zht/super-job/scheduler/internal/service/http/client"
	"github.com/ac-zht/super-job/scheduler/pkg/utils"
	"go.uber.org/zap"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Slack struct {
	repo          repository.SettingRepository
	httpClient    *client.HttpClient
	timeout       time.Duration
	retryTimes    uint8
	retryInterval time.Duration
}

func NewSlackNotify(repo repository.SettingRepository,
	httpClient *client.HttpClient,
	opts ...option.Option[Slack]) Notifiable {
	slack := &Slack{
		repo:          repo,
		httpClient:    httpClient,
		timeout:       time.Second,
		retryTimes:    3,
		retryInterval: time.Second,
	}
	option.Apply[Slack](slack, opts...)
	return slack
}

func WithTimeout(timeout time.Duration) option.Option[Slack] {
	return func(service *Slack) {
		service.timeout = timeout
	}
}

func WithSlackRetryTimes(rt uint8) option.Option[Slack] {
	return func(service *Slack) {
		service.retryTimes = rt
	}
}

func WithSlackRetryInterval(ri time.Duration) option.Option[Slack] {
	return func(service *Slack) {
		service.retryInterval = ri
	}
}

func (s *Slack) Name() string {
	return "slack"
}

func (s *Slack) Send(ctx context.Context, msg Message) {
	slack, err := s.repo.Slack(ctx)
	if err != nil {
		zap.L().Error("#slack#从数据库获取slack配置失败", zap.Error(err))
		return
	}
	if slack.Url == "" {
		zap.L().Error("#slack#Url为空")
		return
	}
	if len(slack.Channels) == 0 {
		zap.L().Error("#slack#channels配置为空")
		return
	}
	channels := s.getActiveSlackChannels(slack, msg)
	msg["content"] = parseNotifyTemplate(slack.Template, msg)
	msg["content"] = html.UnescapeString(msg["content"].(string))
	for _, channel := range channels {
		s.send(msg, slack.Url, channel)
	}
}

func (s *Slack) send(msg Message, slackUrl string, channel string) {
	formatBody := s.format(msg["content"].(string), channel)
	var i uint8
	for {
		s.httpClient.Url = slackUrl
		s.httpClient.Timeout = s.timeout
		resp := s.httpClient.PostJson(formatBody)
		if resp.StatusCode == http.StatusOK {
			break
		}
		i++
		time.Sleep(s.retryInterval)
		if i < s.retryTimes {
			zap.L().Error(fmt.Sprintf("slack#发送消息失败#%s#消息内容-%s", resp.Body, msg["content"]))
			continue
		}
		break
	}
}

func (s *Slack) getActiveSlackChannels(slack domain.Slack, msg Message) []string {
	taskReceiverIds := strings.Split(msg["task_receiver_id"].(string), ",")
	channels := []string{}
	for _, v := range slack.Channels {
		if utils.InStringSlice(taskReceiverIds, strconv.Itoa(int(v.Id))) {
			channels = append(channels, v.Name)
		}
	}
	return channels
}

// 格式化消息内容
func (s *Slack) format(content string, channel string) string {
	content = utils.EscapeJson(content)
	specialChars := []string{"&", "<", ">"}
	replaceChars := []string{"&amp;", "&lt;", "&gt;"}
	content = utils.ReplaceStrings(content, specialChars, replaceChars)

	return fmt.Sprintf(`{"text":"%s","username":"gocron", "channel":"%s"}`, content, channel)
}

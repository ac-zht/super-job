package integration

import (
	"bytes"
	"encoding/json"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/integration/startup"
	"github.com/ac-zht/super-job/admin/internal/repository/dao"
	"github.com/ac-zht/super-job/admin/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ExecutorReq struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Hosts string `json:"hosts"`
}

type WebTestSuite struct {
	suite.Suite
	server *gin.Engine
	db     *gorm.DB
}

func TestGORMArticle(t *testing.T) {
	suite.Run(t, new(WebTestSuite))
}

func (a *WebTestSuite) SetupSuite() {
	a.server = startup.InitWeb()
	a.db = startup.InitTestDB()
}

func (a *WebTestSuite) SetupTest() {
	err := a.db.Exec("TRUNCATE TABLE `tasks`").Error
	assert.NoError(a.T(), err)
	err = a.db.Exec("TRUNCATE TABLE `users`").Error
	assert.NoError(a.T(), err)
	err = a.db.Exec("SET FOREIGN_KEY_CHECKS=0").Error
	err = a.db.Exec("TRUNCATE TABLE `executors`").Error
	assert.NoError(a.T(), err)
	err = a.db.Exec("SET FOREIGN_KEY_CHECKS=1").Error
}

func (a *WebTestSuite) TestExecutorHandler_Save() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    ExecutorReq

		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "新增执行器",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {
				var executor dao.Executor
				a.db.Where("name = ?", "executor-1").First(&executor)
				assert.True(t, executor.Ctime > 0)
				assert.True(t, executor.Utime > 0)
				executor.Utime = 0
				executor.Ctime = 0
				assert.Equal(t, dao.Executor{
					Id:    1,
					Name:  "executor-1",
					Hosts: "10.0.0.100:9000,10.0.0.100:9001",
				}, executor)
			},
			req: ExecutorReq{
				Name:  "executor-1",
				Hosts: "10.0.0.100:9000,10.0.0.100:9001",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Msg:  "新增成功",
				Data: 1,
			},
		},
	}
	for _, tc := range testCases {
		a.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/api/executor/save", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()
			a.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			var result Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

func (a *WebTestSuite) TestTaskHandler_Save() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    web.TaskEditReq

		wantCode   int
		wantResult Result[int64]
	}{
		{
			name: "新增任务",
			before: func(t *testing.T) {
				a.db.Create(&dao.Executor{
					Id:    1,
					Name:  "executor-1",
					Hosts: "10.0.0.100:9001",
					Ctime: time.Now().UnixMilli(),
					Utime: time.Now().UnixMilli(),
				})
				a.db.Create(&dao.User{
					Id:      1,
					Name:    "张三",
					Email:   "2032754457@qq.com",
					Salt:    "qaz123",
					IsAdmin: 1,
					Status:  1,
					Ctime:   time.Now().UnixMilli(),
					Utime:   time.Now().UnixMilli(),
				})
				a.db.Create(&dao.User{
					Id:      2,
					Name:    "李四",
					Email:   "2032754455@qq.com",
					Salt:    "qaz123",
					IsAdmin: 1,
					Status:  1,
					Ctime:   time.Now().UnixMilli(),
					Utime:   time.Now().UnixMilli(),
				})
			},
			after: func(t *testing.T) {
				var task dao.Task
				a.db.Where("name = ?", "task-1").First(&task)
				assert.True(t, task.Ctime > 0)
				assert.True(t, task.Utime > 0)
				task.Utime = 0
				task.Ctime = 0
				expression := "0 * * * * *"
				assert.Equal(t, dao.Task{
					Id:               1,
					Name:             "task-1",
					ExecId:           1,
					Cfg:              "this is a test task",
					Expression:       expression,
					NextTime:         Next(time.Now(), expression).UnixMilli(),
					Status:           1,
					Multi:            0,
					Protocol:         2,
					ExecutorHandler:  "taskHandler",
					Timeout:          10,
					RetryTimes:       3,
					RetryInterval:    1,
					NotifyStatus:     2,
					NotifyType:       1,
					NotifyReceiverId: "8",
				}, task)
			},
			req: web.TaskEditReq{
				Name:             "task-1",
				ExecId:           1,
				Cfg:              "this is a test task",
				Expression:       "0 * * * * *",
				Status:           domain.TaskStatusWaiting,
				Multi:            domain.SingleInstanceRun,
				Protocol:         domain.TaskRPC.ToUint8(),
				ExecutorHandler:  "taskHandler",
				Timeout:          10,
				RetryTimes:       3,
				RetryInterval:    1,
				NotifyStatus:     domain.OverNotification.ToUint8(),
				NotifyType:       domain.EmailNotification.ToUint8(),
				NotifyReceiverId: "8",
			},
			wantCode: 200,
			wantResult: Result[int64]{
				Msg:  "新增成功",
				Data: 1,
			},
		},
	}
	for _, tc := range testCases {
		a.T().Run(tc.name, func(t *testing.T) {
			tc.before(t)
			data, err := json.Marshal(tc.req)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost,
				"/api/task/save", bytes.NewReader(data))
			assert.NoError(t, err)
			req.Header.Set("Content-Type",
				"application/json")
			recorder := httptest.NewRecorder()
			a.server.ServeHTTP(recorder, req)
			code := recorder.Code
			assert.Equal(t, tc.wantCode, code)
			if code != http.StatusOK {
				return
			}
			var result Result[int64]
			err = json.Unmarshal(recorder.Body.Bytes(), &result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
			tc.after(t)
		})
	}
}

func Next(t time.Time, expression string) time.Time {
	expr := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom |
		cron.Month | cron.Dow |
		cron.Descriptor)
	s, _ := expr.Parse(expression)
	return s.Next(t)
}

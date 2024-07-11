package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zc-zht/super-job/admin/internal/integration/startup"
	"github.com/zc-zht/super-job/admin/internal/repository/dao"
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

type JobReq struct {
	Id            int64  `json:"id"`
	ExecId        int64  `json:"exec_id"`
	Name          string `json:"name"`
	Protocol      uint8  `json:"protocol"`
	Cfg           string `json:"cfg"`
	Expression    string `json:"expression"`
	Status        uint8  `json:"status"`
	Multi         uint8  `json:"multi"`
	HttpMethod    uint8  `json:"http_method"`
	Timeout       int64  `json:"timeout"`
	RetryTimes    int64  `json:"retry_times"`
	RetryInterval int64  `json:"retry_interval"`
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
	err := a.db.Exec("TRUNCATE TABLE `jobs`").Error
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
				"/executor/save", bytes.NewReader(data))
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

func (a *WebTestSuite) TestJobHandler_Save() {
	testCases := []struct {
		name   string
		before func(t *testing.T)
		after  func(t *testing.T)
		req    JobReq

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
				var job dao.Job
				a.db.Where("name = ?", "job-1").First(&job)
				assert.True(t, job.Ctime > 0)
				assert.True(t, job.Utime > 0)
				job.Utime = 0
				job.Ctime = 0
				expression := "0 * * * * *"
				assert.Equal(t, dao.Job{
					Id:         1,
					Name:       "job-1",
					ExecId:     1,
					Protocol:   2,
					Cfg:        "this is a test job",
					Expression: expression,
					NextTime:   Next(time.Now(), expression).UnixMilli(),
				}, job)
			},
			req: JobReq{
				ExecId:        1,
				Name:          "job-1",
				Protocol:      2,
				Cfg:           "this is a test job",
				Expression:    "0 * * * * *",
				Status:        1,
				Multi:         1,
				Timeout:       10,
				RetryTimes:    3,
				RetryInterval: 1,
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
				"/job/save", bytes.NewReader(data))
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

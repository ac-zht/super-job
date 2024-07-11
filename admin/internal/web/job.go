package web

import (
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/errs"
	"github.com/zc-zht/super-job/admin/internal/service"
	"github.com/zc-zht/super-job/admin/pkg/ginx"
	"github.com/zc-zht/super-job/admin/pkg/logger"
	"net/http"
	"strconv"
)

type JobHandler struct {
	svc service.JobService
	l   logger.Logger
}

func NewJobHandler(svc service.JobService) *JobHandler {
	return &JobHandler{
		svc: svc,
	}
}

func (h *JobHandler) List(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		return
	}
	offset := (page - 1) * pageSize
	jobs, err := h.svc.List(ctx, offset, pageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.JobInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	data := map[string]interface{}{
		"total": 1,
		"jobs": slice.Map[domain.Job, JobVo](jobs, func(idx int, src domain.Job) JobVo {
			return JobVo{
				Id:            src.Id,
				Executor:      src.Executor.Name,
				Name:          src.Name,
				Protocol:      src.Protocol.ToUint8(),
				Cfg:           src.Cfg,
				Expression:    src.Expression,
				Status:        src.Status,
				Multi:         src.Multi,
				HttpMethod:    src.HttpMethod.ToUint8(),
				Timeout:       src.Timeout,
				RetryTimes:    src.RetryTimes,
				RetryInterval: src.RetryInterval,
				NextTime:      src.NextTime.UnixMilli(),
				Ctime:         src.Ctime,
			}
		}),
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: data,
	})
	return
}

func (h *JobHandler) Save(ctx *gin.Context) {
	type Req struct {
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
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	id, err := h.svc.Save(ctx, domain.Job{
		Id:            req.Id,
		ExecId:        req.ExecId,
		Name:          req.Name,
		Protocol:      domain.JobProtocol(req.Protocol),
		Cfg:           req.Cfg,
		Expression:    req.Expression,
		Status:        req.Status,
		Multi:         req.Multi,
		HttpMethod:    domain.HttpMethod(req.HttpMethod),
		Timeout:       req.Timeout,
		RetryTimes:    req.RetryTimes,
		RetryInterval: req.RetryInterval,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.JobInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg:  "新增成功",
		Data: id,
	})
	return
}

func (h *JobHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.JobInvalidInput,
			Msg:  "参数错误",
		})
		return
	}
	err = h.svc.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.JobInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "删除成功",
	})
	return
}

func (h *JobHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/api/job")
	ug.GET("", h.List)
	ug.POST("/save", h.Save)
	ug.POST("/delete/:id", h.Delete)
}

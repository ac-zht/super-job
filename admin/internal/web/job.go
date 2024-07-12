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
				Id:              src.Id,
				Executor:        src.Executor.Name,
				Name:            src.Name,
				Protocol:        src.Protocol.ToUint8(),
				Cfg:             src.Cfg,
				Expression:      src.Expression,
				Status:          src.Status,
				Multi:           src.Multi,
				HttpMethod:      src.HttpMethod.ToUint8(),
				ExecutorHandler: src.ExecutorHandler,
				Command:         src.Command,
				Timeout:         src.Timeout,
				RetryTimes:      src.RetryTimes,
				RetryInterval:   src.RetryInterval,
				NextTime:        src.NextTime.UnixMilli(),
				Ctime:           src.Ctime,
			}
		}),
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: data,
	})
	return
}

func (h *JobHandler) Save(ctx *gin.Context) {
	var req JobEditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	id, err := h.svc.Save(ctx, domain.Job{
		Id:               req.Id,
		ExecId:           req.ExecId,
		Name:             req.Name,
		Cfg:              req.Cfg,
		Expression:       req.Expression,
		Status:           req.Status,
		Multi:            req.Multi,
		Protocol:         domain.JobProtocol(req.Protocol),
		HttpMethod:       domain.HttpMethod(req.HttpMethod),
		ExecutorHandler:  req.ExecutorHandler,
		Command:          req.Command,
		Timeout:          req.Timeout,
		RetryTimes:       req.RetryTimes,
		RetryInterval:    req.RetryInterval,
		NotifyStatus:     domain.NotifyStatus(req.NotifyStatus),
		NotifyType:       domain.NotifyType(req.NotifyType),
		NotifyReceiverId: req.NotifyReceiverId,
		NotifyKeyword:    req.NotifyKeyword,
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

func (h *JobHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.JobInvalidInput,
			Msg:  "参数错误",
		})
		return
	}
	job, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.JobInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: JobDetail{
			Id:               job.Id,
			ExecId:           job.ExecId,
			Name:             job.Name,
			Cfg:              job.Cfg,
			Expression:       job.Expression,
			Protocol:         job.Protocol.ToUint8(),
			HttpMethod:       job.HttpMethod.ToUint8(),
			Status:           job.Status,
			Multi:            job.Multi,
			ExecutorHandler:  job.ExecutorHandler,
			Command:          job.Command,
			Timeout:          job.Timeout,
			RetryTimes:       job.RetryTimes,
			RetryInterval:    job.RetryInterval,
			NotifyStatus:     job.NotifyStatus.ToUint8(),
			NotifyType:       job.NotifyType.ToUint8(),
			NotifyReceiverId: job.NotifyReceiverId,
			NotifyKeyword:    job.NotifyKeyword,
		},
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
	ug.GET("/:id", h.Detail)
	ug.POST("/save", h.Save)
	ug.POST("/delete/:id", h.Delete)
}

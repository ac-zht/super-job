package web

import (
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/errs"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/pkg/ginx"
	"github.com/ac-zht/super-job/admin/pkg/logger"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	svc service.TaskService
	l   logger.Logger
}

func NewTaskHandler(svc service.TaskService) *TaskHandler {
	return &TaskHandler{
		svc: svc,
	}
}

func (h *TaskHandler) List(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		return
	}
	offset := (page - 1) * pageSize
	tasks, err := h.svc.List(ctx, offset, pageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.TaskInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	data := map[string]interface{}{
		"total": 1,
		"tasks": slice.Map[domain.Task, TaskVo](tasks, func(idx int, src domain.Task) TaskVo {
			return TaskVo{
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

func (h *TaskHandler) Save(ctx *gin.Context) {
	var req TaskEditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	id, err := h.svc.Save(ctx, domain.Task{
		Id:               req.Id,
		ExecId:           req.ExecId,
		Name:             req.Name,
		Cfg:              req.Cfg,
		Expression:       req.Expression,
		Status:           req.Status,
		Multi:            req.Multi,
		Protocol:         domain.TaskProtocol(req.Protocol),
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
			Code: errs.TaskInternalServerError,
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

func (h *TaskHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.TaskInvalidInput,
			Msg:  "参数错误",
		})
		return
	}
	task, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.TaskInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: TaskDetail{
			Id:               task.Id,
			ExecId:           task.ExecId,
			Name:             task.Name,
			Cfg:              task.Cfg,
			Expression:       task.Expression,
			Protocol:         task.Protocol.ToUint8(),
			HttpMethod:       task.HttpMethod.ToUint8(),
			Status:           task.Status,
			Multi:            task.Multi,
			ExecutorHandler:  task.ExecutorHandler,
			Command:          task.Command,
			Timeout:          task.Timeout,
			RetryTimes:       task.RetryTimes,
			RetryInterval:    task.RetryInterval,
			NotifyStatus:     task.NotifyStatus.ToUint8(),
			NotifyType:       task.NotifyType.ToUint8(),
			NotifyReceiverId: task.NotifyReceiverId,
			NotifyKeyword:    task.NotifyKeyword,
		},
	})
	return
}

func (h *TaskHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.TaskInvalidInput,
			Msg:  "参数错误",
		})
		return
	}
	err = h.svc.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.TaskInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "删除成功",
	})
	return
}

func (h *TaskHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/api/task")
	ug.GET("", h.List)
	ug.GET("/:id", h.Detail)
	ug.POST("/save", h.Save)
	ug.POST("/delete/:id", h.Delete)
}

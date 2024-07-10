package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zc-zht/super-job/admin/internal/domain"
	"github.com/zc-zht/super-job/admin/internal/errs"
	"github.com/zc-zht/super-job/admin/internal/service"
	"github.com/zc-zht/super-job/admin/pkg/ginx"
	"github.com/zc-zht/super-job/admin/pkg/logger"
	"net/http"
	"strings"
)

type ExecutorHandler struct {
	svc service.ExecutorService
	l   logger.Logger
}

func NewExecutorHandler(svc service.ExecutorService) *ExecutorHandler {
	return &ExecutorHandler{
		svc: svc,
	}
}

func (h *ExecutorHandler) List(ctx *gin.Context) {
	type Req struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	jobs, err := h.svc.List(ctx, req.Offset, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.JobInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: jobs,
	})
	return
}

func (h *ExecutorHandler) Save(ctx *gin.Context) {
	type Req struct {
		Id    int64  `json:"id"`
		Name  string `json:"name"`
		Hosts string `json:"hosts"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Name == "" || req.Hosts == "" {
		ctx.JSON(http.StatusOK, ginx.Result{Code: 1, Msg: "未输入必填项"})
		return
	}
	id, err := h.svc.Save(ctx, domain.Executor{
		Id:    req.Id,
		Name:  req.Name,
		Hosts: strings.Split(req.Hosts, ","),
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

func (h *ExecutorHandler) Delete(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := h.svc.Delete(ctx, req.Id)
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

func (h *ExecutorHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/executor")
	ug.POST("", h.List)
	ug.POST("/save", h.Save)
	ug.POST("/delete", h.Delete)
}

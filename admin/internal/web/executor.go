package web

import (
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/errs"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/pkg/ginx"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ExecutorHandler struct {
	svc service.ExecutorService
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
	executors, err := h.svc.List(ctx, req.Offset, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.ExecutorInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: slice.Map[domain.Executor, ExecutorVo](executors, func(idx int, src domain.Executor) ExecutorVo {
			return ExecutorVo{
				Id:    src.Id,
				Name:  src.Name,
				Hosts: strings.Join(src.Hosts, "_"),
				Ctime: src.Ctime,
				Utime: src.Utime,
			}
		}),
	})
	return
}

func (h *ExecutorHandler) All(ctx *gin.Context) {
	executors, err := h.svc.List(ctx, 0, -1)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.ExecutorInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: slice.Map[domain.Executor, ExecutorBrief](executors, func(idx int, src domain.Executor) ExecutorBrief {
			return ExecutorBrief{
				Id:   src.Id,
				Name: src.Name,
			}
		}),
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
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.ExecutorRequiredNotInput,
			Msg:  "未输入必填项",
		})
		return
	}
	id, err := h.svc.Save(ctx, domain.Executor{
		Id:    req.Id,
		Name:  req.Name,
		Hosts: strings.Split(req.Hosts, ","),
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.ExecutorInternalServerError,
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
			Code: errs.ExecutorInternalServerError,
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
	ug := server.Group("/api/executor")
	ug.GET("", h.List)
	ug.GET("/all", h.All)
	ug.POST("/save", h.Save)
	ug.POST("/delete", h.Delete)
}

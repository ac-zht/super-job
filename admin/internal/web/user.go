package web

import (
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/errs"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/pkg/ginx"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

func (h *UserHandler) List(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil {
		return
	}
	offset := (page - 1) * pageSize
	users, err := h.svc.List(ctx, offset, pageSize)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	data := map[string]interface{}{
		"total": 1,
		"tasks": users,
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: data,
	})
	return
}

func (h *UserHandler) Save(ctx *gin.Context) {
	var req UserEditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	id, err := h.svc.Save(ctx, domain.User{
		Id:       req.Id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		IsAdmin:  req.IsAdmin,
		Status:   req.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInternalServerError,
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

func (h *UserHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "参数错误",
		})
		return
	}
	user, err := h.svc.GetById(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: user,
	})
	return
}

func (h *UserHandler) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "参数错误",
		})
		return
	}
	err = h.svc.Delete(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "删除成功",
	})
	return
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/api/task")
	ug.GET("", h.List)
	ug.GET("/:id", h.Detail)
	ug.POST("/save", h.Save)
	ug.POST("/delete/:id", h.Delete)
}

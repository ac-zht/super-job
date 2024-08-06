package web

import (
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/errs"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/internal/web/jwt"
	"github.com/ac-zht/super-job/admin/pkg/ginx"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	svc service.UserService
	jwt.Handler
}

func NewUserHandler(svc service.UserService, jwtHdl jwt.Handler) *UserHandler {
	return &UserHandler{
		svc:     svc,
		Handler: jwtHdl,
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
	total, err := h.svc.Count(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	data := map[string]interface{}{
		"total": total,
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
	if req.Name == "" || req.Email == "" {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "用户名或邮箱不能为空",
		})
		return
	}
	if req.Id == 0 {
		if req.Password == "" {
			ctx.JSON(http.StatusOK, ginx.Result{
				Code: errs.UserInvalidInput,
				Msg:  "请输入密码",
			})
			return
		}
		if req.Password != req.ConfirmPassword {
			ctx.JSON(http.StatusOK, ginx.Result{
				Code: errs.UserInvalidInput,
				Msg:  "两次密码输入不一致",
			})
			return
		}
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
		if err == service.ErrUserDuplicate {
			ctx.JSON(http.StatusOK, ginx.Result{
				Code: errs.UserDuplicateUsernameOrEmail,
				Msg:  "用户名或邮箱冲突",
			})
			return
		}
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg:  "成功",
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

func (h *UserHandler) Login(ctx *gin.Context) {
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.Username == "" || req.Password == "" {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInvalidInput,
			Msg:  "用户名或密码不能为空",
		})
		return
	}
	user, err := h.svc.ValidateLogin(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  err.Error(),
		})
		return
	}

	err = h.SetLoginToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.UserInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "登录成功",
		Data: LoginResp{
			Uid:      user.Id,
			Username: user.Name,
			IsAdmin:  user.IsAdmin,
		},
	})
	return
}

func (h *UserHandler) Enable(ctx *gin.Context) {

}

func (h *UserHandler) Disable(ctx *gin.Context) {

}

func (h *UserHandler) UpdateMyPassword(ctx *gin.Context) {

}

func (h *UserHandler) UpdatePassword(ctx *gin.Context) {

}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/api/user")
	ug.GET("", h.List)
	ug.GET("/:id", h.Detail)
	ug.POST("/save", h.Save)
	ug.POST("/remove/:id", h.Delete)
	ug.POST("/login", h.Login)
	ug.POST("/enable/:id", h.Enable)
	ug.POST("/disable/:id", h.Disable)
	ug.POST("/editMyPassword", h.UpdateMyPassword)
	ug.POST("/editPassword/:id", h.UpdatePassword)
}

package web

import (
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/errs"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ac-zht/super-job/admin/pkg/ginx"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InstallHandler struct {
	svc service.InstallService
}

func NewInstallHandler(svc service.InstallService) *InstallHandler {
	return &InstallHandler{
		svc: svc,
	}
}

func (h *InstallHandler) Store(ctx *gin.Context) {
	if service.App.Installed {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.InstallOccurred,
		})
		return
	}
	var req domain.Installation
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if req.AdminPassword != req.ConfirmAdminPassword {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.InstallPasswordInconsistent,
		})
		return
	}
	err := h.svc.Store(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.InstallInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "安装成功",
	})
}

func (h *InstallHandler) Status(ctx *gin.Context) {
	installed, _ := h.svc.Status(ctx)
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: installed,
	})
}

func (h *InstallHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("api/install")
	ug.GET("/store", h.Store)
	ug.GET("/status", h.Status)
}

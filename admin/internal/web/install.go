package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zc-zht/super-job/admin/internal/errs"
	"github.com/zc-zht/super-job/admin/internal/service"
	"github.com/zc-zht/super-job/admin/pkg/ginx"
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
	if service.Installed {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.InstallOccurred,
		})
		return
	}
	err := h.svc.Store(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.InstallInternalServerError,
		})
		return
	}
	err = service.CreateInstallLock()
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.InstallInternalServerError,
			Msg:  "创建文件安装锁失败",
		})
		return
	}
	service.Installed = true
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "安装成功",
	})
}

func (h *InstallHandler) Status(ctx *gin.Context) {

}

func (h *InstallHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("api/install")
	ug.GET("/store", h.Store)
	ug.GET("/status", h.Status)
}

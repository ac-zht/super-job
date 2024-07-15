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

func (h *InstallHandler) Store(ctx *gin.Context) {
	err := h.svc.Store(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.InstallInternalServerError,
		})
		return
	}
}

func (h *InstallHandler) Status(ctx *gin.Context) {

}

func (h *InstallHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("api/install")
	ug.GET("/store", h.Store)
	ug.GET("/status", h.Status)
}

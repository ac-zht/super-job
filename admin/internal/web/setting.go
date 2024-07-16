package web

import (
	"github.com/gin-gonic/gin"
	"github.com/zc-zht/super-job/admin/internal/errs"
	"github.com/zc-zht/super-job/admin/internal/service"
	"github.com/zc-zht/super-job/admin/pkg/ginx"
	"net/http"
)

type SettingHandler struct {
	svc service.SettingService
}

func NewSettingHandler(svc service.SettingService) *SettingHandler {
	return &SettingHandler{
		svc: svc,
	}
}

func (h *SettingHandler) Mail(ctx *gin.Context) {
	mail, err := h.svc.Mail(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: mail,
	})
}

func (h *SettingHandler) UpdateMail(ctx *gin.Context) {

}

func (h *SettingHandler) CreateMailUser(ctx *gin.Context) {

}

func (h *SettingHandler) RemoveMailUser(ctx *gin.Context) {

}

func (h *SettingHandler) Slack(ctx *gin.Context) {
	slack, err := h.svc.Slack(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: slack,
	})
}

func (h *SettingHandler) Webhook(ctx *gin.Context) {
	webhook, err := h.svc.WebHook(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
			Msg:  "系统异常",
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Data: webhook,
	})
}

func (h *SettingHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("api/system")
	mail := ug.Group("/mail")
	mail.GET("", h.Mail)
	mail.POST("/update", h.UpdateMail)
	mail.POST("/user", h.CreateMailUser)
	mail.POST("/user/remove/:id", h.RemoveMailUser)
	slack := ug.Group("/slack")
	slack.GET("", h.Slack)
	webhook := ug.Group("/webhook")
	webhook.GET("", h.Webhook)

}

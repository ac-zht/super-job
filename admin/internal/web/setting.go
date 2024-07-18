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
	var req MailEditReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInvalidInput,
		})
		return
	}
	err := h.svc.UpdateMail(ctx, domain.MailServer{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
	}, req.Template)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "邮箱配置成功",
	})
}

func (h *SettingHandler) CreateMailUser(ctx *gin.Context) {
	var mailUser domain.MailUser
	if err := ctx.Bind(&mailUser); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInvalidInput,
		})
		return
	}
	id, err := h.svc.CreateMailUser(ctx, mailUser)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg:  "添加接收用户成功",
		Data: id,
	})
}

func (h *SettingHandler) RemoveMailUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInvalidInput,
		})
		return
	}
	err = h.svc.RemoveMailUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "删除成功",
	})
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

func (h *SettingHandler) UpdateSlack(ctx *gin.Context) {
	var slack domain.Slack
	if err := ctx.Bind(&slack); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInvalidInput,
		})
		return
	}
	err := h.svc.UpdateSlack(ctx, slack)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "slack配置成功",
	})
}

func (h *SettingHandler) CreateChannel(ctx *gin.Context) {
	var channel domain.Channel
	if err := ctx.Bind(&channel); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInvalidInput,
		})
		return
	}
	id, err := h.svc.CreateChannel(ctx, channel)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg:  "添加channel成功",
		Data: id,
	})
}

func (h *SettingHandler) RemoveChannel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInvalidInput,
		})
		return
	}
	err = h.svc.RemoveChannel(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "删除成功",
	})
}

func (h *SettingHandler) Webhook(ctx *gin.Context) {
	webhook, err := h.svc.Webhook(ctx)
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

func (h *SettingHandler) UpdateWebhook(ctx *gin.Context) {
	var webhook domain.Webhook
	if err := ctx.Bind(&webhook); err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInvalidInput,
		})
		return
	}
	err := h.svc.UpdateWebhook(ctx, webhook)
	if err != nil {
		ctx.JSON(http.StatusOK, ginx.Result{
			Code: errs.SettingInternalServerError,
		})
		return
	}
	ctx.JSON(http.StatusOK, ginx.Result{
		Msg: "webhook配置成功",
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
	slack.POST("/update", h.UpdateSlack)
	slack.POST("/channel", h.CreateChannel)
	slack.POST("/channel/remove/:id", h.RemoveChannel)
	webhook := ug.Group("/webhook")
	webhook.GET("", h.Webhook)
	webhook.POST("/update", h.UpdateWebhook)

}

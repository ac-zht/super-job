package middleware

import (
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
)

type LoginJWTMiddlewareBuilder struct {
	publicPaths set.Set[string]
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	s := set.NewMapSet[string](3)
	s.Add("/api")
	s.Add("/api/user/login")
	s.Add("/api/install/status")
	return &LoginJWTMiddlewareBuilder{
		publicPaths: s,
	}
}

func (m *LoginJWTMiddlewareBuilder) RestoreToken(ctx *gin.Context) error {
	authToken := ctx.Request.Header.Get("Auth-Token")
	if authToken == "" {
		return nil
	}
}

func (m *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !service.App.Installed {
			return
		}

		if m.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
	}
}

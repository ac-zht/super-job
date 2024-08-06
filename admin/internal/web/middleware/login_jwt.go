package middleware

import (
	"errors"
	"github.com/ac-zht/super-job/admin/internal/service"
	mjwt "github.com/ac-zht/super-job/admin/internal/web/jwt"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type LoginJWTMiddlewareBuilder struct {
	publicPaths set.Set[string]
	mjwt.Handler
}

func NewLoginJWTMiddlewareBuilder(handler mjwt.Handler) *LoginJWTMiddlewareBuilder {
	s := set.NewMapSet[string](3)
	s.Add("/api")
	s.Add("/api/user/login")
	s.Add("/api/install/status")
	return &LoginJWTMiddlewareBuilder{
		publicPaths: s,
		Handler:     handler,
	}
}

func (m *LoginJWTMiddlewareBuilder) ParseJwtToken(ctx *gin.Context) (mjwt.UserClaims, error) {
	uc := mjwt.UserClaims{}
	authToken := ctx.Request.Header.Get("Auth-Token")
	if authToken == "" {
		return uc, nil
	}
	token, err := jwt.ParseWithClaims(authToken, uc, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.App.Setting.AccessTokenKey), nil
	})
	if err != nil || !token.Valid {
		return uc, errors.New("invalid token")
	}
	expireTime, err := uc.GetExpirationTime()
	if err != nil {
		return uc, err
	}
	if expireTime.Before(time.Now()) {
		return uc, errors.New("token expired")
	}
	return uc, nil
}

func (m *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !service.App.Installed {
			return
		}
		if m.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
		uc, err := m.ParseJwtToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if m.CheckSession(ctx, uc.Ssid) != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("user", uc)
	}
}

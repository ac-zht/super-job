package middleware

import (
	"errors"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/ecodeclub/ekit/set"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.App.Setting.AuthSecret), nil
	})
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid claims")
	}
	ctx.Set("uid", int(claims["uid"].(float64)))
	ctx.Set("username", claims["username"])
	ctx.Set("is_admin", int(claims["is_admin"].(float64)))
	return nil
}

func (m *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !service.App.Installed {
			return
		}
		m.RestoreToken(ctx)
		if m.publicPaths.Exist(ctx.Request.URL.Path) {
			return
		}
	}
}

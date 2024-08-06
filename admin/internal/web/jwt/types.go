package jwt

import (
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	ClearToken(ctx *gin.Context) error
	SetLoginToken(ctx *gin.Context, user domain.User) error
	SetJWTToken(ctx *gin.Context, ssid string, user domain.User) error
	CheckSession(ctx *gin.Context, ssid string) error
	ExtractTokenString(ctx *gin.Context) string
}

type UserClaims struct {
	Id       int64
	Ssid     string
	Username string
	IsAdmin  uint8
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	Id   int64
	Ssid string
	jwt.RegisteredClaims
}

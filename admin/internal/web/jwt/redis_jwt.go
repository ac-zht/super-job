package jwt

import (
	"errors"
	"fmt"
	"github.com/ac-zht/super-job/admin/internal/domain"
	"github.com/ac-zht/super-job/admin/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"time"
)

type RedisHandler struct {
	cmd redis.Cmdable
	exp time.Duration
}

// ClearToken 标记用户Ssid的请求退出态让其token作废，防止退出后再用token操作，存储时长为refresh_token时长
func (r *RedisHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-refresh-token", "")
	uc := ctx.MustGet("user").(UserClaims)
	return r.cmd.Set(ctx, r.key(uc.Ssid), "", r.exp).Err()
}

func (r *RedisHandler) SetLoginToken(ctx *gin.Context, user domain.User) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, ssid, user)
	if err != nil {
		return err
	}
	err = r.setRefreshToken(ctx, ssid, user.Id)
	return err
}

func (r *RedisHandler) SetJWTToken(ctx *gin.Context, ssid string, user domain.User) error {
	uc := UserClaims{
		Id:       user.Id,
		Ssid:     ssid,
		Username: user.Name,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "super-job",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(domain.TokenDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	tokenStr, err := token.SignedString([]byte(service.App.Setting.AccessTokenKey))
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (r *RedisHandler) setRefreshToken(ctx *gin.Context, ssid string, uid int64) error {
	rc := RefreshClaims{
		Id:   uid,
		Ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "super-job",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(domain.RefreshTokenDuration)),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rc)
	refreshTokenStr, err := refreshToken.SignedString([]byte(service.App.Setting.RefreshTokenKey))
	if err != nil {
		return err
	}
	ctx.Header("x-refresh-token", refreshTokenStr)
	return nil
}

func (r *RedisHandler) CheckSession(ctx *gin.Context, ssid string) error {
	logout, err := r.cmd.Exists(ctx, r.key(ssid)).Result()
	if err != nil {
		return err
	}
	if logout > 0 {
		return errors.New("user exited")
	}
	return nil
}

func (r *RedisHandler) ExtractTokenString(ctx *gin.Context) string {
	authCode := ctx.GetHeader("Authorization")
	if authCode == "" {
		return ""
	}
	authSegments := strings.SplitN(authCode, " ", 2)
	if len(authSegments) != 2 {
		return ""
	}
	return authSegments[1]
}

func (r *RedisHandler) key(ssid string) string {
	return fmt.Sprintf("user:Ssid:%s", ssid)
}

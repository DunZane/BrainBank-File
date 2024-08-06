package middleware

import (
	"context"
	"github.com/dunzane/brainbank-file/api/internal/config"
	"github.com/dunzane/brainbank-file/api/internal/utils"
	"net/http"
	"strings"
)

type JwtMiddleware struct {
	Config config.Config
}

func NewJwtMiddleware(c config.Config) *JwtMiddleware {
	return &JwtMiddleware{c}
}

func (m *JwtMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取token
		authHeader := r.Header.Get("Authorization")
		token := strings.Replace(authHeader, "Bearer ", "", 1)

		//// 判断是否携带token
		//if token == "" {
		//	// 直接响应Token为空
		//	http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		//	return
		//}

		// 解析token
		claims, err := utils.ParseToken(m.Config, token)
		if err != nil {
			// 直接响应Token不合法
			http.Error(w, "Invalid or malformed token", http.StatusUnauthorized)
			return
		}

		//// 校验token是否过期
		//if time.Now().Unix() > claims.ExpiresAt.Unix() {
		//	// 直接响应Token已经过期需要重新尝试登录
		//	http.Error(w, "Token has expired, please log in again", http.StatusUnauthorized)
		//	return
		//}

		// 存储用户信息在context
		reqCtx := r.Context()
		ctx := context.WithValue(reqCtx, "userId", claims.UserId)
		ctx = context.WithValue(ctx, "email", claims.Email)
		newReq := r.WithContext(ctx)
		next(w, newReq)
	}
}

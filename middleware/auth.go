package middleware

import (
	"context"
	"llm-ops/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("your-secret-key"))

func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// log.Printf("=== Auth Check for %s %s ===", r.Method, r.URL.Path)

		// 打印所有请求头以便调试
		// log.Printf("Request Headers:")
		// for name, values := range r.Header {
		// 	log.Printf("  %s: %v", name, values)
		// }

		token := r.Header.Get("Authorization")
		if token == "" {
			log.Printf("No Authorization header found")
			utils.ErrorResponse(w, http.StatusUnauthorized, "未提供认证信息")
			return
		}

		// 移除 "Bearer " 前缀
		token = strings.TrimPrefix(token, "Bearer ")

		claims, err := utils.ValidateToken(token)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			utils.ErrorResponse(w, http.StatusUnauthorized, "无效的认证信息")
			return
		}

		// log.Printf("Token validated successfully for user: %s (role: %s)", claims.Username, claims.Role)

		// 将用户信息添加到请求上下文
		ctx := context.WithValue(r.Context(), "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// 添加管理员权限检查中间件
func AdminRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := utils.GetUserFromContext(r.Context())
		if !ok || claims.Role != "admin" {
			log.Printf("Admin access denied for user: %+v", claims)
			utils.ErrorResponse(w, http.StatusForbidden, "需要管理员权限")
			return
		}
		next.ServeHTTP(w, r)
	}
}

// 添加Nebula权限检查中间件 - 只允许运维和管理员访问
func NebulaRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, ok := utils.GetUserFromContext(r.Context())
		if !ok {
			log.Printf("No user found in context")
			utils.ErrorResponse(w, http.StatusUnauthorized, "未获取到用户信息")
			return
		}

		if claims.Role != "admin" && claims.Role != "operator" {
			log.Printf("Nebula access denied for user: %s with role: %s", claims.Username, claims.Role)
			utils.ErrorResponse(w, http.StatusForbidden, "需要运维或管理员权限")
			return
		}
		next.ServeHTTP(w, r)
	}
}

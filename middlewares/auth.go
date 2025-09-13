package middlewares

import (
	"context"
	"net/http"
	"os"
	"project/utils"
	"strings"
)

const (
	ContextUserID    = "userID"
	ContextUserEmail = "userEmail"
	ContextUserRole  = "userRole"
)

// AuthMiddleware — проверяет наличие и валидность accessToken (JWT).
// Берём из Authorization: Bearer <token> или из cookie["accessToken"].
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		// 1. Пробуем Authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Если нет — пробуем cookie
		if tokenStr == "" {
			cookie, err := r.Cookie("accessToken")
			if err == nil {
				tokenStr = cookie.Value
			}
		}

		if tokenStr == "" {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		// 3. Парсим JWT
		claims, err := utils.ParseToken(tokenStr, os.Getenv("JWT_SECRET_KEY"))
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// 4. Извлекаем данные
		email, _ := claims["email"].(string)
		role, _ := claims["role"].(string)
		id, _ := claims["id"].(string) // если в токен пишешь userID

		// 5. Кладём в контекст
		ctx := context.WithValue(r.Context(), ContextUserID, id)
		ctx = context.WithValue(ctx, ContextUserEmail, email)
		ctx = context.WithValue(ctx, ContextUserRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RefreshMiddleware — проверяет наличие refreshToken.
func RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("refreshToken")
		if err != nil || token.Value == "" {
			http.Error(w, "Unauthorized: missing refresh token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RoleMiddleware — универсальный чекер по ролям.
// Пример: router.Handle("/admin", RoleMiddleware("ADMIN")(adminHandler))
func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextUserRole).(string)
			if !ok || role == "" {
				http.Error(w, "Forbidden: no role", http.StatusForbidden)
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Forbidden: insufficient rights", http.StatusForbidden)
		})
	}
}

// Helpers

// GetUserEmail возвращает email из контекста
func GetUserEmail(ctx context.Context) string {
	if v, ok := ctx.Value(ContextUserEmail).(string); ok {
		return v
	}
	return ""
}

// GetUserRole возвращает роль из контекста
func GetUserRole(ctx context.Context) string {
	if v, ok := ctx.Value(ContextUserRole).(string); ok {
		return v
	}
	return ""
}

// GetUserID возвращает userID из контекста
func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(ContextUserID).(string); ok {
		return v
	}
	return ""
}

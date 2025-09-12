package middlewares

import (
	"context"
	// "fmt"
	"net/http"
	"os"
	"project/utils"
	// "strings"
)

const (
	ContextUserEmail = "userEmail"
	ContextUserRole  = "userRole"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := r.Cookie("accessToken")
		// fmt.Println("Access Token:", tokenStr)
		if err != nil || tokenStr.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(tokenStr.Value, os.Getenv("JWT_SECRET_KEY"))
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		email := claims["email"].(string)
		role := claims["role"].(string)

		ctx := context.WithValue(r.Context(), ContextUserEmail, email)
		ctx = context.WithValue(ctx, ContextUserRole, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RefreshMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := r.Cookie("refreshToken")
		if err != nil || refreshToken.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func OnlyAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("accessToken")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(cookie.Value, os.Getenv("JWT_SECRET_KEY"))
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		roleClaim, ok := claims["role"].(string)
		// fmt.Println("Role Claim:", roleClaim)
		if !ok || roleClaim != "admin" {
			http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

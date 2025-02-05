package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Beeram12/college-appointment-system/pkg/utils"
)

type contextKey string

const UserContextKey contextKey = "userClaims"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorisation header", http.StatusUnauthorized)
			return
		}
		// Extracting token from the handler
		tokenParts := strings.Split(authHeader, "")
		if len(tokenParts) == 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid authorzation format", http.StatusUnauthorized)
			return
		}
		tokenString := tokenParts[1]
		claims, err := utils.ValidatingTokens(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

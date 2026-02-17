package middleware

import (
	"context"
	"incident-report-app/models"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, `{"error":"missing or invalid authorization header"}`, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// Reject any algorithm that isn't HMAC (e.g. "none" attacks)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret()), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		// JWT numbers decode as float64 in Go
		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, `{"error":"invalid token payload"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), models.UserIDKey, int64(userID))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jwtSecret() string {
	if s := os.Getenv("JWT_SECRET"); s != "" {
		return s
	}
	return "dev-secret-change-in-production"
}
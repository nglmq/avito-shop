package middleware

import (
	"context"
	ujwt "github.com/nglmq/avito-shop/internal/utils/jwt"
	"log/slog"
	"net/http"
)

const (
	ContextUserID = "user"
)

func CheckAuthMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Authorization token is required", http.StatusUnauthorized)
				return
			}

			tokenString = tokenString[7:]

			userID, err := ujwt.GetUserID(tokenString)
			if err != nil {
				logger.Error("Invalid auth token",
					slog.String("path", r.URL.Path),
					slog.String("error", err.Error()))
				http.Error(w, "Invalid auth token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

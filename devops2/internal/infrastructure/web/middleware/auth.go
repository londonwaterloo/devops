package middleware

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"bmstu_devops_lab_2/internal/domain/service"
	"context"
	"net/http"
)

// AuthMiddleware проверяет наличие валидной сессии
func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем cookie сессии
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Валидируем сессию
			user, err := authService.ValidateSession(cookie.Value)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Добавляем сессию в контекст
			ctx := context.WithValue(r.Context(), "session", &entity.Session{
				ID:       cookie.Value,
				UserID:   user.ID,
				Username: user.Username,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware добавляет сессию в контекст если она существует, но не блокирует запрос
func OptionalAuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Пытаемся получить cookie сессии
			cookie, err := r.Cookie("session_id")
			if err == nil {
				// Валидируем сессию если cookie есть
				user, err := authService.ValidateSession(cookie.Value)
				if err == nil {
					// Добавляем сессию в контекст
					ctx = context.WithValue(ctx, "session", &entity.Session{
						ID:       cookie.Value,
						UserID:   user.ID,
						Username: user.Username,
					})
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

package handler

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"bmstu_devops_lab_2/internal/domain/service"
	"encoding/json"
	"net/http"
)

// AuthHandler обрабатывает запросы аутентификации
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler создаёт новый экземпляр AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResponse представляет ответ на регистрацию
type RegisterResponse struct {
	UserID   int64  `json:"user_id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse представляет ответ на вход
type LoginResponse struct {
	SessionID string `json:"session_id"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
}

// ErrorResponse представляет ответ с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// Register обрабатывает запрос на регистрацию
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := ah.authService.Register(req.Username, req.Password)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Устанавливаем cookie сессии
	sessionID, err := ah.authService.Login(req.Username, req.Password)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setSessionCookie(w, sessionID.ID)

	resp := RegisterResponse{
		UserID:   user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// Login обрабатывает запрос на вход
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, err := ah.authService.Login(req.Username, req.Password)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	setSessionCookie(w, session.ID)

	resp := LoginResponse{
		SessionID: session.ID,
		UserID:    session.UserID,
		Username:  session.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Logout обрабатывает запрос на выход
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем сессию из контекста
	session, ok := r.Context().Value("session").(*entity.Session)
	if !ok {
		sendErrorResponse(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	if err := ah.authService.Logout(session.ID); err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Удаляем cookie
	clearSessionCookie(w)

	w.WriteHeader(http.StatusNoContent)
}

// GetCurrentUser возвращает информацию о текущем пользователе
func (ah *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем сессию из контекста
	session, ok := r.Context().Value("session").(*entity.Session)
	if !ok {
		sendErrorResponse(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	resp := map[string]interface{}{
		"user_id":   session.UserID,
		"username":  session.Username,
		"session_id": session.ID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// sendErrorResponse отправляет JSON-ответ с ошибкой
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// setSessionCookie устанавливает cookie сессии
func setSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   24 * 60 * 60, // 24 часа
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// clearSessionCookie удаляет cookie сессии
func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

package web

import (
	"bmstu_devops_lab_2/internal/domain/service"
	"bmstu_devops_lab_2/internal/handler"
	"bmstu_devops_lab_2/internal/infrastructure/web/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

// Router настраивает маршруты HTTP
type Router struct {
	*mux.Router
	fileHandler  *handler.FileHandler
	authHandler  *handler.AuthHandler
	templatesDir string
	authService *service.AuthService
}

// NewRouter создаёт новый роутер
func NewRouter(fileHandler *handler.FileHandler, authHandler *handler.AuthHandler, templatesDir string, authService *service.AuthService) *Router {
	r := &Router{
		Router:       mux.NewRouter(),
		fileHandler:  fileHandler,
		authHandler:  authHandler,
		templatesDir: templatesDir,
		authService: authService,
	}
	r.setupRoutes()
	return r
}

// ServeHTTP реализует интерфейс http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

// setupRoutes настраивает маршруты
func (r *Router) setupRoutes() {
	// Auth routes (не требуют аутентификации)
	r.Router.HandleFunc("/api/auth/register", r.authHandler.Register).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/auth/login", r.authHandler.Login).Methods(http.MethodPost)

	// Auth routes (требуют аутентификации)
	authAPI := r.Router.PathPrefix("/api/auth").Subrouter()
	authAPI.Use(middleware.AuthMiddleware(r.authService))
	authAPI.HandleFunc("/logout", r.authHandler.Logout).Methods(http.MethodPost)

	// Auth routes с опциональной аутентификацией (для проверки статуса)
	r.Router.Handle("/api/auth/me", middleware.OptionalAuthMiddleware(r.authService)(http.HandlerFunc(r.authHandler.GetCurrentUser))).Methods(http.MethodGet)

	// File routes (публичные)
	r.Router.HandleFunc("/api/files", r.fileHandler.ListFiles).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/files/{id}/download", r.fileHandler.DownloadFile).Methods(http.MethodGet)

	// File routes (требуют аутентификации)
	api := r.Router.PathPrefix("/api/files").Subrouter()
	api.Use(middleware.AuthMiddleware(r.authService))
	api.HandleFunc("/upload", r.fileHandler.UploadFile).Methods(http.MethodPost)
	api.HandleFunc("/{id}", r.fileHandler.DeleteFile).Methods(http.MethodDelete)

	// Web routes (HTML)
	r.Router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, r.templatesDir+"/base.html")
	}).Methods(http.MethodGet)
}

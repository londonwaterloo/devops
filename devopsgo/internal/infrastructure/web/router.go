package web

import (
	"bmstu_devops_lab_1/internal/handler"
	"net/http"

	"github.com/gorilla/mux"
)

// Router настраивает маршруты HTTP
type Router struct {
	*mux.Router
	fileHandler   *handler.FileHandler
	templatesDir  string
}

// NewRouter создаёт новый роутер
func NewRouter(fileHandler *handler.FileHandler, templatesDir string) *Router {
	r := &Router{
		Router:       mux.NewRouter(),
		fileHandler:   fileHandler,
		templatesDir:  templatesDir,
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
	// API routes
	r.Router.HandleFunc("/api/files", r.fileHandler.ListFiles).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/files/upload", r.fileHandler.UploadFile).Methods(http.MethodPost)
	r.Router.HandleFunc("/api/files/{id}/download", r.fileHandler.DownloadFile).Methods(http.MethodGet)
	r.Router.HandleFunc("/api/files/{id}", r.fileHandler.DeleteFile).Methods(http.MethodDelete)

	// Web routes (HTML)
	r.Router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, r.templatesDir+"/base.html")
	}).Methods(http.MethodGet)
}

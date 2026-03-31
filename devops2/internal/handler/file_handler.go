package handler

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"bmstu_devops_lab_2/internal/infrastructure/storage"
	"bmstu_devops_lab_2/internal/usecase"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// FileHandler обрабатывает HTTP-запросы
type FileHandler struct {
	useCase *usecase.FileUseCase
}

// NewFileHandler создаёт новый экземпляр FileHandler
func NewFileHandler(useCase *usecase.FileUseCase) *FileHandler {
	return &FileHandler{useCase: useCase}
}

// ListFiles возвращает список файлов в формате JSON
func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	files, err := h.useCase.GetAllFiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// UploadFile обрабатывает загрузку файла
func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Проверяем аутентификацию
	session, ok := r.Context().Value("session").(*entity.Session)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Ограничиваем размер загружаемого файла до 32MB
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)

	file, header, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			http.Error(w, "file is required", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	fileEntity := &entity.File{
		ID:        storage.GenerateID(),
		Name:      header.Filename,
		Size:      int64(len(content)),
		Content:   content,
		CreatedAt: time.Now(),
		Uploader:  session.Username, // сохраняем имя пользователя
	}

	if err := h.useCase.UploadFile(fileEntity); err != nil {
		if err.Error() == "file already exists" {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fileEntity)
}

// DeleteFile обрабатывает удаление файла
func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	// Проверяем аутентификацию
	session, ok := r.Context().Value("session").(*entity.Session)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	// Получаем информацию о файле
	file, err := h.useCase.GetFile(id)
	if err != nil {
		if err.Error() == "file not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Проверяем, что файл принадлежит текущему пользователю
	if file.Uploader != session.Username {
		http.Error(w, "Forbidden: you can only delete your own files", http.StatusForbidden)
		return
	}

	if err := h.useCase.DeleteFile(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DownloadFile обрабатывает скачивание файла
func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	file, err := h.useCase.GetFile(id)
	if err != nil {
		if err.Error() == "file not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Disposition", `attachment; filename="`+file.Name+`"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(file.Content)
}

package repository

import "bmstu_devops_lab_2/internal/domain/entity"

// FileRepository определяет интерфейс для работы с файлами
type FileRepository interface {
	// Save сохраняет файл в хранилище
	Save(file *entity.File) error

	// FindAll возвращает список всех файлов
	FindAll() ([]*entity.File, error)

	// FindByID возвращает файл по ID
	FindByID(id string) (*entity.File, error)

	// ReadContent читает содержимое файла по ID
	ReadContent(id string) ([]byte, string, error)

	// Delete удаляет файл по ID
	Delete(id string) error

	// Exists проверяет существование файла
	Exists(name string) bool

	// ExistsForUploader проверяет существование файла для конкретного пользователя
	ExistsForUploader(name, uploader string) bool
}

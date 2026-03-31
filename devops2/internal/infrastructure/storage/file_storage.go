package storage

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const metadataFile = ".metadata.json"

// FileStorage реализует FileRepository на основе файловой системы
type FileStorage struct {
	uploadDir    string
	files        map[string]*entity.File
	mu           sync.RWMutex
	metadataPath string
}

// fileMetadata структура для хранения метаданных в файле
type fileMetadata struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
}

// NewFileStorage создаёт новый экземпляр FileStorage
func NewFileStorage(uploadDir string) *FileStorage {
	return &FileStorage{
		uploadDir:    uploadDir,
		files:        make(map[string]*entity.File),
		metadataPath: filepath.Join(uploadDir, metadataFile),
	}
}

// LoadMetadata загружает метаданные из файла
func (fs *FileStorage) LoadMetadata() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	data, err := os.ReadFile(fs.metadataPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Файл метаданных не существует - это нормально для первого запуска
			return nil
		}
		return err
	}

	var metadata []fileMetadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return err
	}

	for _, m := range metadata {
		createdAt, err := time.Parse(time.RFC3339, m.CreatedAt)
		if err != nil {
			continue // Пропускаем записи с некорректной датой
		}

		// Проверяем, что файл существует на диске
		filePath := filepath.Join(fs.uploadDir, m.Name)
		if _, err := os.Stat(filePath); err == nil {
			fs.files[m.ID] = &entity.File{
				ID:        m.ID,
				Name:      m.Name,
				Size:      m.Size,
				CreatedAt: createdAt,
			}
		}
	}

	return nil
}

// saveMetadata сохраняет метаданные в файл (без блокировки, должен вызываться под Lock)
func (fs *FileStorage) saveMetadata() error {
	var metadata []fileMetadata
	for _, file := range fs.files {
		metadata = append(metadata, fileMetadata{
			ID:        file.ID,
			Name:      file.Name,
			Size:      file.Size,
			CreatedAt: file.CreatedAt.Format(time.RFC3339),
		})
	}

	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fs.metadataPath, data, 0644)
}

// SaveMetadata сохраняет метаданные в файл
func (fs *FileStorage) SaveMetadata() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	return fs.saveMetadata()
}

// Save сохраняет файл в хранилище
func (fs *FileStorage) Save(file *entity.File) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	filePath := filepath.Join(fs.uploadDir, file.Name)
	if err := os.WriteFile(filePath, file.Content, 0644); err != nil {
		return err
	}

	fs.files[file.ID] = file
	return fs.saveMetadata()
}

// FindAll возвращает список всех файлов
func (fs *FileStorage) FindAll() ([]*entity.File, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	result := make([]*entity.File, 0, len(fs.files))
	for _, file := range fs.files {
		result = append(result, file)
	}
	return result, nil
}

// FindByID возвращает файл по ID
func (fs *FileStorage) FindByID(id string) (*entity.File, error) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	file, ok := fs.files[id]
	if !ok {
		return nil, errors.New("file not found")
	}
	return file, nil
}

// Delete удаляет файл по ID
func (fs *FileStorage) Delete(id string) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	file, ok := fs.files[id]
	if !ok {
		return errors.New("file not found")
	}

	filePath := filepath.Join(fs.uploadDir, file.Name)
	if err := os.Remove(filePath); err != nil {
		return err
	}

	delete(fs.files, id)
	return fs.saveMetadata()
}

// Exists проверяет существование файла
func (fs *FileStorage) Exists(name string) bool {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	for _, file := range fs.files {
		if file.Name == name {
			return true
		}
	}
	return false
}

// ExistsForUploader проверяет существование файла для конкретного пользователя
func (fs *FileStorage) ExistsForUploader(name, uploader string) bool {
	// В FileStorage нет понятия uploader, возвращаем false
	return false
}

// ReadContent читает содержимое файла по ID
func (fs *FileStorage) ReadContent(id string) ([]byte, string, error) {
	fs.mu.RLock()
	file, ok := fs.files[id]
	fs.mu.RUnlock()

	if !ok {
		return nil, "", errors.New("file not found")
	}

	filePath := filepath.Join(fs.uploadDir, file.Name)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	return content, file.Name, nil
}

// InitUploadDir создаёт директорию для загрузки файлов
func (fs *FileStorage) InitUploadDir() error {
	return os.MkdirAll(fs.uploadDir, 0755)
}

// CreateEntity создаёт новую сущность File
func (fs *FileStorage) CreateEntity(name string, size int64) *entity.File {
	return &entity.File{
		ID:        GenerateID(),
		Name:      name,
		Size:      size,
		CreatedAt: time.Now(),
	}
}

func GenerateID() string {
	return time.Now().Format("20060102150405.999")
}

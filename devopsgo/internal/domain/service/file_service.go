package service

// FileService определяет интерфейс сервисов домена
type FileService interface {
	// ValidateName проверяет валидность имени файла
	ValidateName(name string) error

	// ValidateSize проверяет допустимый размер файла
	ValidateSize(size int64) error
}

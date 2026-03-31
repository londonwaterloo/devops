package usecase

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"bmstu_devops_lab_2/internal/domain/repository"
	"bmstu_devops_lab_2/internal/domain/service"
	"errors"
)

// FileUseCase содержит бизнес-логику работы с файлами
type FileUseCase struct {
	repo repository.FileRepository
	svc  service.FileService
}

// NewFileUseCase создаёт новый экземпляр FileUseCase
func NewFileUseCase(repo repository.FileRepository, svc service.FileService) *FileUseCase {
	return &FileUseCase{
		repo: repo,
		svc:  svc,
	}
}

// GetAllFiles возвращает список всех файлов
func (uc *FileUseCase) GetAllFiles() ([]*entity.File, error) {
	return uc.repo.FindAll()
}

// UploadFile загружает новый файл
func (uc *FileUseCase) UploadFile(file *entity.File) error {
	if err := uc.svc.ValidateName(file.Name); err != nil {
		return err
	}
	if err := uc.svc.ValidateSize(file.Size); err != nil {
		return err
	}
	// Проверяем, что файл с таким именем уже не существует для этого пользователя
	if uc.repo.ExistsForUploader(file.Name, file.Uploader) {
		return errors.New("file already exists")
	}
	return uc.repo.Save(file)
}

// DeleteFile удаляет файл по ID
func (uc *FileUseCase) DeleteFile(id string) error {
	if _, err := uc.repo.FindByID(id); err != nil {
		return errors.New("file not found")
	}
	return uc.repo.Delete(id)
}

// GetFile возвращает файл с содержимым по ID
func (uc *FileUseCase) GetFile(id string) (*entity.File, error) {
	file, err := uc.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("file not found")
	}

	content, _, err := uc.repo.ReadContent(id)
	if err != nil {
		return nil, err
	}

	file.Content = content
	return file, nil
}

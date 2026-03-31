package repository

import "bmstu_devops_lab_2/internal/domain/entity"

// SessionRepository определяет интерфейс для работы с сессиями
type SessionRepository interface {
	// Create создаёт новую сессию
	Create(session *entity.Session) error

	// FindByID находит сессию по ID
	FindByID(id string) (*entity.Session, error)

	// Delete удаляет сессию по ID
	Delete(id string) error

	// DeleteByUserID удаляет все сессии пользователя
	DeleteByUserID(userID int64) error
}

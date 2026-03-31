package repository

import "bmstu_devops_lab_2/internal/domain/entity"

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	// Create создаёт нового пользователя
	Create(user *entity.User) error

	// FindByUsername находит пользователя по имени
	FindByUsername(username string) (*entity.User, error)

	// FindByID находит пользователя по ID
	FindByID(id int64) (*entity.User, error)

	// ExistsByUsername проверяет существование пользователя по имени
	ExistsByUsername(username string) (bool, error)
}

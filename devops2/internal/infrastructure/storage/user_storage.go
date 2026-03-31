package storage

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// PostgresUserStorage реализует UserRepository на основе PostgreSQL
type PostgresUserStorage struct {
	db *sql.DB
}

// NewPostgresUserStorage создаёт новый экземпляр PostgresUserStorage
func NewPostgresUserStorage(db *sql.DB) *PostgresUserStorage {
	return &PostgresUserStorage{db: db}
}

// Create создаёт нового пользователя
func (pus *PostgresUserStorage) Create(user *entity.User) error {
	query := `
		INSERT INTO users (username, password, created_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := pus.db.QueryRow(query, user.Username, user.Password, user.CreatedAt).
		Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.New("username already exists")
		}
		return err
	}
	return nil
}

// FindByUsername находит пользователя по имени
func (pus *PostgresUserStorage) FindByUsername(username string) (*entity.User, error) {
	query := `
		SELECT id, username, password, created_at
		FROM users
		WHERE username = $1
	`
	var user entity.User
	err := pus.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// FindByID находит пользователя по ID
func (pus *PostgresUserStorage) FindByID(id int64) (*entity.User, error) {
	query := `
		SELECT id, username, password, created_at
		FROM users
		WHERE id = $1
	`
	var user entity.User
	err := pus.db.QueryRow(query, id).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// ExistsByUsername проверяет существование пользователя по имени
func (pus *PostgresUserStorage) ExistsByUsername(username string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE username = $1
		)
	`
	var exists bool
	err := pus.db.QueryRow(query, username).Scan(&exists)
	return exists, err
}

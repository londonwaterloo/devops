package service

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"bmstu_devops_lab_2/internal/domain/repository"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService содержит бизнес-логику аутентификации и авторизации
type AuthService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

// NewAuthService создаёт новый экземпляр AuthService
func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

// Register регистрирует нового пользователя
func (as *AuthService) Register(username, password string) (*entity.User, error) {
	// Проверяем длину имени пользователя
	if len(username) < 3 || len(username) > 50 {
		return nil, errors.New("username must be between 3 and 50 characters")
	}

	// Проверяем длину пароля
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Проверяем, что пользователь уже не существует
	exists, err := as.userRepo.ExistsByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаём пользователя
	user := &entity.User{
		Username:  username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	if err := as.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Не возвращаем пароль
	user.Password = ""
	return user, nil
}

// Login выполняет вход пользователя и возвращает ID сессии
func (as *AuthService) Login(username, password string) (*entity.Session, error) {
	// Находим пользователя
	user, err := as.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Создаём сессию
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	now := time.Now()
	session := &entity.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Username:  user.Username,
		CreatedAt: now,
		ExpiresAt: now.Add(24 * time.Hour), // сессия действительна 24 часа
	}

	if err := as.sessionRepo.Create(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// Logout выполняет выход пользователя по ID сессии
func (as *AuthService) Logout(sessionID string) error {
	return as.sessionRepo.Delete(sessionID)
}

// ValidateSession проверяет валидность сессии и возвращает данные пользователя
func (as *AuthService) ValidateSession(sessionID string) (*entity.User, error) {
	session, err := as.sessionRepo.FindByID(sessionID)
	if err != nil {
		return nil, errors.New("invalid session")
	}

	// Проверяем срок действия сессии
	if time.Now().After(session.ExpiresAt) {
		// Удаляем просроченную сессию
		_ = as.sessionRepo.Delete(sessionID)
		return nil, errors.New("session expired")
	}

	return &entity.User{
		ID:       session.UserID,
		Username: session.Username,
	}, nil
}

// LogoutAll выполняет выход пользователя из всех сессий
func (as *AuthService) LogoutAll(userID int64) error {
	return as.sessionRepo.DeleteByUserID(userID)
}

// generateSessionID генерирует уникальный ID сессии
func generateSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

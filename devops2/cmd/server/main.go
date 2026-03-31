package main

import (
	"bmstu_devops_lab_2/internal/config"
	"bmstu_devops_lab_2/internal/domain/service"
	"bmstu_devops_lab_2/internal/handler"
	"bmstu_devops_lab_2/internal/infrastructure/logger"
	"bmstu_devops_lab_2/internal/infrastructure/storage"
	"bmstu_devops_lab_2/internal/infrastructure/web"
	"bmstu_devops_lab_2/internal/infrastructure/web/middleware"
	"bmstu_devops_lab_2/internal/usecase"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	logger := logger.NewStdLogger()
	logger.Info("Starting file server...")

	// Инициализация PostgreSQL
	db, err := sql.Open("postgres", cfg.Postgres.DSNIPv4())
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer db.Close()

	// Проверяем подключение к PostgreSQL
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping PostgreSQL: %v", err)
	}

	// Инициализация хранилища пользователей
	userStorage := storage.NewPostgresUserStorage(db)

	// Инициализация Redis
	sessionStorage, err := storage.NewRedisSessionStorage(cfg.Redis.Addr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer sessionStorage.Close()

	// Инициализация хранилища файлов с MongoDB
	fileStorage, err := storage.NewMongoStorage(
		cfg.MongoDB.URI,
		cfg.MongoDB.Database,
		cfg.MongoDB.Collection,
		cfg.Storage.UploadDir,
	)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer fileStorage.Close()

	if err := fileStorage.InitUploadDir(); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Инициализация сервисов домена
	fileService := service.NewFileServiceImpl(cfg.Storage.MaxFileSize)
	authService := service.NewAuthService(userStorage, sessionStorage)

	// Инициализация usecase
	fileUseCase := usecase.NewFileUseCase(fileStorage, fileService)

	// Инициализация обработчиков
	fileHandler := handler.NewFileHandler(fileUseCase)
	authHandler := handler.NewAuthHandler(authService)

	// Инициализация роутера
	router := middleware.RecoveryMiddleware(
		middleware.LoggingMiddleware(
			web.NewRouter(fileHandler, authHandler, cfg.Templates.Dir, authService),
		),
	)

	// Запуск сервера
	addr := cfg.Address()
	logger.Info(fmt.Sprintf("Server listening on %s", addr))

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

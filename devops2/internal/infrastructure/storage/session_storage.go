package storage

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisSessionStorage реализует SessionRepository на основе Redis
type RedisSessionStorage struct {
	client *redis.Client
}

// NewRedisSessionStorage создаёт новый экземпляр RedisSessionStorage
func NewRedisSessionStorage(addr string) (*RedisSessionStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisSessionStorage{client: client}, nil
}

// Close закрывает соединение с Redis
func (rss *RedisSessionStorage) Close() error {
	return rss.client.Close()
}

// Create создаёт новую сессию
func (rss *RedisSessionStorage) Create(session *entity.Session) error {
	ctx := context.Background()

	key := "session:" + session.ID

	// Сохраняем данные сессии как hash
	err := rss.client.HSet(ctx, key, map[string]interface{}{
		"user_id":    session.UserID,
		"username":   session.Username,
		"created_at": session.CreatedAt.Unix(),
		"expires_at": session.ExpiresAt.Unix(),
	}).Err()
	if err != nil {
		return err
	}

	// Устанавливаем время жизни ключа (TTL) до expires_at
	ttl := time.Until(session.ExpiresAt)
	if ttl > 0 {
		return rss.client.Expire(ctx, key, ttl).Err()
	}
	return nil
}

// FindByID находит сессию по ID
func (rss *RedisSessionStorage) FindByID(id string) (*entity.Session, error) {
	ctx := context.Background()
	key := "session:" + id

	data, err := rss.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		return nil, errors.New("session not found")
	}

	userID, _ := strconv.ParseInt(data["user_id"], 10, 64)
	createdAt, _ := strconv.ParseInt(data["created_at"], 10, 64)
	expiresAt, _ := strconv.ParseInt(data["expires_at"], 10, 64)

	return &entity.Session{
		ID:        id,
		UserID:    userID,
		Username:  data["username"],
		CreatedAt: time.Unix(createdAt, 0),
		ExpiresAt: time.Unix(expiresAt, 0),
	}, nil
}

// Delete удаляет сессию по ID
func (rss *RedisSessionStorage) Delete(id string) error {
	ctx := context.Background()
	key := "session:" + id
	return rss.client.Del(ctx, key).Err()
}

// DeleteByUserID удаляет все сессии пользователя
func (rss *RedisSessionStorage) DeleteByUserID(userID int64) error {
	ctx := context.Background()

	// Находим все ключи сессий для этого пользователя
	var cursor uint64
	var keys []string

	for {
		var batch []string
		var err error

		batch, cursor, err = rss.client.Scan(ctx, cursor, "session:*", 100).Result()
		if err != nil {
			return err
		}

		// Проверяем каждую сессию
		for _, key := range batch {
			sessionID := key[8:] // убираем префикс "session:"
			session, err := rss.FindByID(sessionID)
			if err == nil && session.UserID == userID {
				keys = append(keys, key)
			}
		}

		if cursor == 0 {
			break
		}
	}

	// Удаляем найденные сессии
	if len(keys) > 0 {
		return rss.client.Del(ctx, keys...).Err()
	}

	return nil
}

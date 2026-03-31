package storage

import (
	"bmstu_devops_lab_2/internal/domain/entity"
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoFileMetadata структура для хранения метаданных в MongoDB
type MongoFileMetadata struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	Size      int64     `bson:"size"`
	CreatedAt time.Time `bson:"created_at"`
	Uploader  string    `bson:"uploader"` // никнейм пользователя, загрузившего файл
}

// MongoStorage реализует FileRepository с использованием MongoDB для метаданных
type MongoStorage struct {
	uploadDir   string
	client      *mongo.Client
	collection  *mongo.Collection
}

// NewMongoStorage создаёт новый экземпляр MongoStorage
func NewMongoStorage(uri, database, collection, uploadDir string) (*MongoStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	coll := client.Database(database).Collection(collection)

	return &MongoStorage{
		uploadDir:  uploadDir,
		client:     client,
		collection: coll,
	}, nil
}

// Close закрывает соединение с MongoDB
func (ms *MongoStorage) Close() error {
	return ms.client.Disconnect(context.Background())
}

// Save сохраняет файл в хранилище (метаданные в MongoDB, содержимое на диске)
func (ms *MongoStorage) Save(file *entity.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Сохраняем файл на диск
	filePath := filepath.Join(ms.uploadDir, file.Name)
	if err := os.WriteFile(filePath, file.Content, 0644); err != nil {
		return err
	}

	// Сохраняем метаданные в MongoDB
	metadata := MongoFileMetadata{
		ID:        file.ID,
		Name:      file.Name,
		Size:      file.Size,
		CreatedAt: file.CreatedAt,
		Uploader:  file.Uploader,
	}

	_, err := ms.collection.InsertOne(ctx, metadata)
	return err
}

// FindAll возвращает список всех файлов
func (ms *MongoStorage) FindAll() ([]*entity.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ms.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var metadatas []MongoFileMetadata
	if err := cursor.All(ctx, &metadatas); err != nil {
		return nil, err
	}

	result := make([]*entity.File, 0, len(metadatas))
	for _, m := range metadatas {
		result = append(result, &entity.File{
			ID:        m.ID,
			Name:      m.Name,
			Size:      m.Size,
			CreatedAt: m.CreatedAt,
			Uploader:  m.Uploader,
		})
	}

	return result, nil
}

// FindByID возвращает файл по ID
func (ms *MongoStorage) FindByID(id string) (*entity.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var metadata MongoFileMetadata
	err := ms.collection.FindOne(ctx, bson.M{"id": id}).Decode(&metadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	return &entity.File{
		ID:        metadata.ID,
		Name:      metadata.Name,
		Size:      metadata.Size,
		CreatedAt: metadata.CreatedAt,
		Uploader:  metadata.Uploader,
	}, nil
}

// ReadContent читает содержимое файла по ID
func (ms *MongoStorage) ReadContent(id string) ([]byte, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var metadata MongoFileMetadata
	err := ms.collection.FindOne(ctx, bson.M{"id": id}).Decode(&metadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, "", errors.New("file not found")
		}
		return nil, "", err
	}

	filePath := filepath.Join(ms.uploadDir, metadata.Name)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}

	return content, metadata.Name, nil
}

// Delete удаляет файл по ID
func (ms *MongoStorage) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Сначала получаем метаданные
	var metadata MongoFileMetadata
	err := ms.collection.FindOne(ctx, bson.M{"id": id}).Decode(&metadata)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("file not found")
		}
		return err
	}

	// Удаляем файл с диска
	filePath := filepath.Join(ms.uploadDir, metadata.Name)
	if err := os.Remove(filePath); err != nil {
		return err
	}

	// Удаляем метаданные из MongoDB
	_, err = ms.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

// Exists проверяет существование файла
func (ms *MongoStorage) Exists(name string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, _ := ms.collection.CountDocuments(ctx, bson.M{"name": name})
	return count > 0
}

// ExistsForUploader проверяет существование файла для конкретного пользователя
func (ms *MongoStorage) ExistsForUploader(name, uploader string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, _ := ms.collection.CountDocuments(ctx, bson.M{"name": name, "uploader": uploader})
	return count > 0
}

// InitUploadDir создаёт директорию для загрузки файлов
func (ms *MongoStorage) InitUploadDir() error {
	return os.MkdirAll(ms.uploadDir, 0755)
}

// CreateEntity создаёт новую сущность File
func (ms *MongoStorage) CreateEntity(name string, size int64) *entity.File {
	return &entity.File{
		ID:        generateID(),
		Name:      name,
		Size:      size,
		CreatedAt: time.Now(),
	}
}

// generateID генерирует уникальный ID на основе ObjectId MongoDB и текущего времени
func generateID() string {
	return primitive.NewObjectID().Hex()
}

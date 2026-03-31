package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config содержит конфигурацию приложения
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Storage   StorageConfig   `yaml:"storage"`
	Templates TemplatesConfig `yaml:"templates"`
	MongoDB   MongoDBConfig   `yaml:"mongodb"`
	Postgres  PostgresConfig  `yaml:"postgres"`
	Redis     RedisConfig     `yaml:"redis"`
}

// ServerConfig конфигурация HTTP-сервера
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// StorageConfig конфигурация хранилища файлов
type StorageConfig struct {
	UploadDir  string `yaml:"upload_dir"`
	MaxFileSize int64  `yaml:"max_file_size"`
}

// TemplatesConfig конфигурация шаблонов
type TemplatesConfig struct {
	Dir string `yaml:"dir"`
}

// MongoDBConfig конфигурация MongoDB
type MongoDBConfig struct {
	URI        string `yaml:"uri"`
	Database   string `yaml:"database"`
	Collection string `yaml:"collection"`
}

// PostgresConfig конфигурация PostgreSQL
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

// DSN возвращает строку подключения к PostgreSQL
func (p *PostgresConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.SSLMode)
}

// DSNIPv4 возвращает строку подключения к PostgreSQL с принудительным IPv4
func (p *PostgresConfig) DSNIPv4() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName)
}

// RedisConfig конфигурация Redis
type RedisConfig struct {
	Addr string `yaml:"addr"`
}

// Load загружает конфигурацию из файла
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

// Address возвращает адрес сервера
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

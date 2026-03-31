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

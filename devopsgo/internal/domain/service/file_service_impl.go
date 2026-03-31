package service

import (
	"errors"
	"fmt"
	"strings"
)

type FileServiceImpl struct {
	maxFileSize int64
}

func NewFileServiceImpl(maxFileSize int64) *FileServiceImpl {
	return &FileServiceImpl{maxFileSize: maxFileSize}
}

func (s *FileServiceImpl) ValidateName(name string) error {
	if name == "" {
		return errors.New("file name cannot be empty")
	}
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("file name contains invalid characters: %s", name)
	}
	return nil
}

func (s *FileServiceImpl) ValidateSize(size int64) error {
	if size <= 0 {
		return errors.New("file size must be positive")
	}
	if s.maxFileSize > 0 && size > s.maxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.maxFileSize)
	}
	return nil
}

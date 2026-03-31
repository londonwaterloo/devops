package utils

import "os"

// EnsureDir создаёт директорию, если она не существует
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// Exists проверяет существование пути
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

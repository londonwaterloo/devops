package logger

import (
	"log"
	"os"
)

// Logger интерфейс для логирования
type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

// StdLogger реализация логгера через стандартный log
type StdLogger struct {
	info  *log.Logger
	error *log.Logger
	debug *log.Logger
}

// NewStdLogger создаёт новый логгер
func NewStdLogger() *StdLogger {
	return &StdLogger{
		info:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
	}
}

// Info логирует информационное сообщение
func (l *StdLogger) Info(msg string) {
	l.info.Println(msg)
}

// Error логирует ошибку
func (l *StdLogger) Error(msg string) {
	l.error.Println(msg)
}

// Debug логирует отладочное сообщение
func (l *StdLogger) Debug(msg string) {
	l.debug.Println(msg)
}

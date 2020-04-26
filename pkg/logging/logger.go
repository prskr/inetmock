//go:generate mockgen -source=logger.go -destination=./../../internal/mock/logging/logger_mock.go -package=logging_mock
package logging

import "go.uber.org/zap"

type Logger interface {
	Named(s string) Logger
	With(fields ...zap.Field) Logger
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

type logger struct {
	underlyingLogger *zap.Logger
}

func NewLogger(underlyingLogger *zap.Logger) *logger {
	return &logger{underlyingLogger: underlyingLogger}
}

func (l logger) Named(s string) Logger {
	return NewLogger(l.underlyingLogger.Named(s))
}

func (l logger) With(fields ...zap.Field) Logger {
	return NewLogger(l.underlyingLogger.With(fields...))
}

func (l logger) Debug(msg string, fields ...zap.Field) {
	l.underlyingLogger.Debug(msg, fields...)
}

func (l logger) Info(msg string, fields ...zap.Field) {
	l.underlyingLogger.Info(msg, fields...)
}

func (l logger) Warn(msg string, fields ...zap.Field) {
	l.underlyingLogger.Warn(msg, fields...)
}

func (l logger) Error(msg string, fields ...zap.Field) {
	l.underlyingLogger.Error(msg, fields...)
}

func (l logger) Panic(msg string, fields ...zap.Field) {
	l.underlyingLogger.Panic(msg, fields...)
}

func (l logger) Fatal(msg string, fields ...zap.Field) {
	l.underlyingLogger.Fatal(msg, fields...)
}

func (l logger) Sync() error {
	return l.underlyingLogger.Sync()
}

package logging

import (
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggingConfig = zap.NewProductionConfig()

type (
	LoggingOption interface {
		Apply(cfg *zap.Config)
	}
	LoggingOptionFunc func(cfg *zap.Config)
)

func (f LoggingOptionFunc) Apply(cfg *zap.Config) {
	f(cfg)
}

func WithLevel(level zap.AtomicLevel) LoggingOption {
	return LoggingOptionFunc(func(cfg *zap.Config) {
		cfg.Level = level
	})
}

func WithDevelopment(developmentLogging bool) LoggingOption {
	return LoggingOptionFunc(func(cfg *zap.Config) {
		cfg.Development = developmentLogging
	})
}

func WithEncoding(encoding string) LoggingOption {
	return LoggingOptionFunc(func(cfg *zap.Config) {
		cfg.Encoding = encoding
	})
}

func WithInitialFields(initialFields map[string]any) LoggingOption {
	return LoggingOptionFunc(func(cfg *zap.Config) {
		cfg.InitialFields = initialFields
	})
}

func ConfigureLogging(opts ...LoggingOption) error {
	for idx := range opts {
		opts[idx].Apply(&loggingConfig)
	}
	if defaultLogger, err := loggingConfig.Build(zap.AddCallerSkip(1)); err != nil {
		return err
	} else {
		zap.ReplaceGlobals(defaultLogger)
	}

	return nil
}

func ParseLevel(levelString string) zap.AtomicLevel {
	switch strings.ToLower(levelString) {
	case "debug":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
}

func CreateLogger() Logger {
	return NewLogger(zap.L())
}

func CreateTestLogger(tb testing.TB) Logger {
	tb.Helper()
	logger := &testLogger{
		testRunning: make(chan struct{}),
		tb:          tb,
		encoder:     zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
	}

	tb.Cleanup(func() {
		close(logger.testRunning)
	})

	return logger
}

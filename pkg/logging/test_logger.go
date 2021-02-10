package logging

import (
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testLogger struct {
	name    string
	fields  []zap.Field
	tb      testing.TB
	encoder zapcore.Encoder
}

func (t testLogger) Named(s string) Logger {
	return testLogger{
		encoder: t.encoder,
		name:    s,
		tb:      t.tb,
		fields:  t.fields,
	}
}

func (t testLogger) With(fields ...zap.Field) Logger {
	return &testLogger{
		encoder: t.encoder,
		name:    t.name,
		fields:  append(t.fields, fields...),
		tb:      t.tb,
	}
}

func (t testLogger) Debug(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.DebugLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(2)),
	}, append(t.fields, fields...))

	if err == nil {
		t.tb.Log(buf.String())
	}
}

func (t testLogger) Info(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.InfoLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(2)),
	}, append(t.fields, fields...))

	if err == nil {
		t.tb.Log(buf.String())
	}
}

func (t testLogger) Warn(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.WarnLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(2)),
	}, append(t.fields, fields...))

	if err == nil {
		t.tb.Log(buf.String())
	}
}

func (t testLogger) Error(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.ErrorLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(2)),
		Stack:      string(debug.Stack()),
	}, append(t.fields, fields...))

	if err == nil {
		t.tb.Log(buf.String())
	}
}

func (t testLogger) Panic(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.PanicLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(2)),
		Stack:      string(debug.Stack()),
	}, append(t.fields, fields...))

	if err == nil {
		t.tb.Error(buf.String())
	}
}

func (t testLogger) Fatal(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.FatalLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(2)),
		Stack:      string(debug.Stack()),
	}, append(t.fields, fields...))

	if err == nil {
		t.tb.Error(buf.String())
	}
}

func (t testLogger) Sync() error {
	return nil
}

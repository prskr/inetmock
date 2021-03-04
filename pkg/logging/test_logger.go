package logging

import (
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const wrapperCallFramesCount = 2

type testLogger struct {
	name        string
	fields      []zap.Field
	tb          testing.TB
	encoder     zapcore.Encoder
	testRunning chan struct{}
}

func (t *testLogger) testFinished() bool {
	select {
	case _, more := <-t.testRunning:
		return !more
	default:
		return false
	}
}

func (t *testLogger) Named(s string) Logger {
	return &testLogger{
		encoder:     t.encoder,
		name:        s,
		tb:          t.tb,
		fields:      t.fields,
		testRunning: t.testRunning,
	}
}

func (t *testLogger) With(fields ...zap.Field) Logger {
	return &testLogger{
		encoder:     t.encoder,
		name:        t.name,
		fields:      append(t.fields, fields...),
		tb:          t.tb,
		testRunning: t.testRunning,
	}
}

func (t *testLogger) Debug(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.DebugLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(wrapperCallFramesCount)),
	}, append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Info(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.InfoLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(wrapperCallFramesCount)),
	}, append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Warn(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.WarnLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(wrapperCallFramesCount)),
	}, append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Error(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.ErrorLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(wrapperCallFramesCount)),
		Stack:      string(debug.Stack()),
	}, append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Panic(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.PanicLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(wrapperCallFramesCount)),
		Stack:      string(debug.Stack()),
	}, append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Error(buf.String())
	}
}

func (t *testLogger) Fatal(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(zapcore.Entry{
		Level:      zapcore.FatalLevel,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(wrapperCallFramesCount)),
		Stack:      string(debug.Stack()),
	}, append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Error(buf.String())
	}
}

func (t *testLogger) Sync() error {
	return nil
}

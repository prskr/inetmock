package logging

import (
	"runtime"
	"runtime/debug"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const wrapperCallFramesCount = 4

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
	buf, err := t.encoder.EncodeEntry(t.entry(msg, zapcore.DebugLevel, false), append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Info(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(t.entry(msg, zapcore.InfoLevel, false), append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Warn(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(t.entry(msg, zapcore.WarnLevel, false), append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Error(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(t.entry(msg, zapcore.ErrorLevel, true), append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Panic(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(t.entry(msg, zapcore.PanicLevel, true), append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) Fatal(msg string, fields ...zap.Field) {
	t.tb.Helper()
	buf, err := t.encoder.EncodeEntry(t.entry(msg, zapcore.FatalLevel, true), append(t.fields, fields...))

	if err == nil && !t.testFinished() {
		t.tb.Log(buf.String())
	}
}

func (t *testLogger) entry(msg string, lvl zapcore.Level, includeStack bool) zapcore.Entry {
	e := zapcore.Entry{
		Level:      lvl,
		Time:       time.Now(),
		LoggerName: t.name,
		Message:    msg,
		Caller:     zapcore.NewEntryCaller(runtime.Caller(wrapperCallFramesCount)),
	}
	if includeStack {
		e.Stack = string(debug.Stack())
	}
	return e
}

func (t *testLogger) Sync() error {
	return nil
}

func (t *testLogger) ZapLogger() *zap.Logger {
	panic("implement me")
}

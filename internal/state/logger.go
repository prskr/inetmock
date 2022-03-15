package state

import (
	"testing"

	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
)

var _ badger.Logger = (*Logger)(nil)

type Logger struct {
	*zap.SugaredLogger
}

func (l Logger) Errorf(s string, i ...any) {
	l.SugaredLogger.Errorf(s, i...)
}

func (l Logger) Warningf(s string, i ...any) {
	l.SugaredLogger.Warnf(s, i...)
}

func (l Logger) Infof(s string, i ...any) {
	l.SugaredLogger.Infof(s, i...)
}

func (l Logger) Debugf(s string, i ...any) {
	l.SugaredLogger.Debugf(s, i...)
}

var _ badger.Logger = (*TestLogger)(nil)

type TestLogger struct {
	testing.TB
}

func (l TestLogger) Errorf(s string, i ...any) {
	l.TB.Errorf(s, i...)
}

func (l TestLogger) Warningf(s string, i ...any) {
	l.TB.Logf(s, i...)
}

func (l TestLogger) Infof(s string, i ...any) {
	l.TB.Logf(s, i...)
}

func (l TestLogger) Debugf(s string, i ...any) {
	l.TB.Logf(s, i...)
}

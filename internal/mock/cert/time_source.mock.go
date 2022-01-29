// Code generated by MockGen. DO NOT EDIT.
// Source: time_source.go

// Package cert_mock is a generated GoMock package.
package cert_mock

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockTimeSource is a mock of TimeSource interface.
type MockTimeSource struct {
	ctrl     *gomock.Controller
	recorder *MockTimeSourceMockRecorder
}

// MockTimeSourceMockRecorder is the mock recorder for MockTimeSource.
type MockTimeSourceMockRecorder struct {
	mock *MockTimeSource
}

// NewMockTimeSource creates a new mock instance.
func NewMockTimeSource(ctrl *gomock.Controller) *MockTimeSource {
	mock := &MockTimeSource{ctrl: ctrl}
	mock.recorder = &MockTimeSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimeSource) EXPECT() *MockTimeSourceMockRecorder {
	return m.recorder
}

// UTCNow mocks base method.
func (m *MockTimeSource) UTCNow() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UTCNow")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// UTCNow indicates an expected call of UTCNow.
func (mr *MockTimeSourceMockRecorder) UTCNow() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UTCNow", reflect.TypeOf((*MockTimeSource)(nil).UTCNow))
}

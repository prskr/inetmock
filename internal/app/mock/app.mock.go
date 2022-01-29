// Code generated by MockGen. DO NOT EDIT.
// Source: app.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cobra "github.com/spf13/cobra"

	logging "gitlab.com/inetmock/inetmock/pkg/logging"
)

// MockApp is a mock of App interface.
type MockApp struct {
	ctrl     *gomock.Controller
	recorder *MockAppMockRecorder
}

// MockAppMockRecorder is the mock recorder for MockApp.
type MockAppMockRecorder struct {
	mock *MockApp
}

// NewMockApp creates a new mock instance.
func NewMockApp(ctrl *gomock.Controller) *MockApp {
	mock := &MockApp{ctrl: ctrl}
	mock.recorder = &MockAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockApp) EXPECT() *MockAppMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockApp) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockAppMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockApp)(nil).Context))
}

// Logger mocks base method.
func (m *MockApp) Logger() logging.Logger {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logger")
	ret0, _ := ret[0].(logging.Logger)
	return ret0
}

// Logger indicates an expected call of Logger.
func (mr *MockAppMockRecorder) Logger() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logger", reflect.TypeOf((*MockApp)(nil).Logger))
}

// MustRun mocks base method.
func (m *MockApp) MustRun() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "MustRun")
}

// MustRun indicates an expected call of MustRun.
func (mr *MockAppMockRecorder) MustRun() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MustRun", reflect.TypeOf((*MockApp)(nil).MustRun))
}

// RootCommand mocks base method.
func (m *MockApp) RootCommand() *cobra.Command {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RootCommand")
	ret0, _ := ret[0].(*cobra.Command)
	return ret0
}

// RootCommand indicates an expected call of RootCommand.
func (mr *MockAppMockRecorder) RootCommand() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RootCommand", reflect.TypeOf((*MockApp)(nil).RootCommand))
}

// Shutdown mocks base method.
func (m *MockApp) Shutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockAppMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockApp)(nil).Shutdown))
}

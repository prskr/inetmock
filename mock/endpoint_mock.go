// Code generated by MockGen. DO NOT EDIT.
// Source: endpoint.go

// Package mock is a generated GoMock package.
package mock

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockEndpoint is a mock of Endpoint interface
type MockEndpoint struct {
	ctrl     *gomock.Controller
	recorder *MockEndpointMockRecorder
}

// MockEndpointMockRecorder is the mock recorder for MockEndpoint
type MockEndpointMockRecorder struct {
	mock *MockEndpoint
}

// NewMockEndpoint creates a new mock instance
func NewMockEndpoint(ctrl *gomock.Controller) *MockEndpoint {
	mock := &MockEndpoint{ctrl: ctrl}
	mock.recorder = &MockEndpointMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockEndpoint) EXPECT() *MockEndpointMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockEndpoint) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockEndpointMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockEndpoint)(nil).Start))
}

// Shutdown mocks base method
func (m *MockEndpoint) Shutdown() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown")
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown
func (mr *MockEndpointMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockEndpoint)(nil).Shutdown))
}

// Name mocks base method
func (m *MockEndpoint) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockEndpointMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockEndpoint)(nil).Name))
}

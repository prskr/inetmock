// Code generated by MockGen. DO NOT EDIT.
// Source: registration.go

// Package endpoint_mock is a generated GoMock package.
package endpoint_mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	endpoint "inetmock.icb4dc0.de/inetmock/internal/endpoint"
)

// MockHandlerRegistry is a mock of HandlerRegistry interface.
type MockHandlerRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerRegistryMockRecorder
}

// MockHandlerRegistryMockRecorder is the mock recorder for MockHandlerRegistry.
type MockHandlerRegistryMockRecorder struct {
	mock *MockHandlerRegistry
}

// NewMockHandlerRegistry creates a new mock instance.
func NewMockHandlerRegistry(ctrl *gomock.Controller) *MockHandlerRegistry {
	mock := &MockHandlerRegistry{ctrl: ctrl}
	mock.recorder = &MockHandlerRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandlerRegistry) EXPECT() *MockHandlerRegistryMockRecorder {
	return m.recorder
}

// AvailableHandlers mocks base method.
func (m *MockHandlerRegistry) AvailableHandlers() []endpoint.HandlerReference {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AvailableHandlers")
	ret0, _ := ret[0].([]endpoint.HandlerReference)
	return ret0
}

// AvailableHandlers indicates an expected call of AvailableHandlers.
func (mr *MockHandlerRegistryMockRecorder) AvailableHandlers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailableHandlers", reflect.TypeOf((*MockHandlerRegistry)(nil).AvailableHandlers))
}

// HandlerForName mocks base method.
func (m *MockHandlerRegistry) HandlerForName(handlerRef endpoint.HandlerReference) (endpoint.ProtocolHandler, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HandlerForName", handlerRef)
	ret0, _ := ret[0].(endpoint.ProtocolHandler)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// HandlerForName indicates an expected call of HandlerForName.
func (mr *MockHandlerRegistryMockRecorder) HandlerForName(handlerRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandlerForName", reflect.TypeOf((*MockHandlerRegistry)(nil).HandlerForName), handlerRef)
}

// RegisterHandler mocks base method.
func (m *MockHandlerRegistry) RegisterHandler(handlerRef endpoint.HandlerReference, handlerProvider endpoint.HandlerProvider) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterHandler", handlerRef, handlerProvider)
}

// RegisterHandler indicates an expected call of RegisterHandler.
func (mr *MockHandlerRegistryMockRecorder) RegisterHandler(handlerRef, handlerProvider interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterHandler", reflect.TypeOf((*MockHandlerRegistry)(nil).RegisterHandler), handlerRef, handlerProvider)
}

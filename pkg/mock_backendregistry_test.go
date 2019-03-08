// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/asecurityteam/transportd/pkg (interfaces: BackendRegistry)

package transportd

import (
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBackendRegistry is a mock of BackendRegistry interface
type MockBackendRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockBackendRegistryMockRecorder
}

// MockBackendRegistryMockRecorder is the mock recorder for MockBackendRegistry
type MockBackendRegistryMockRecorder struct {
	mock *MockBackendRegistry
}

// NewMockBackendRegistry creates a new mock instance
func NewMockBackendRegistry(ctrl *gomock.Controller) *MockBackendRegistry {
	mock := &MockBackendRegistry{ctrl: ctrl}
	mock.recorder = &MockBackendRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBackendRegistry) EXPECT() *MockBackendRegistryMockRecorder {
	return m.recorder
}

// Load mocks base method
func (m *MockBackendRegistry) Load(arg0 context.Context, arg1 string) http.RoundTripper {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", arg0, arg1)
	ret0, _ := ret[0].(http.RoundTripper)
	return ret0
}

// Load indicates an expected call of Load
func (mr *MockBackendRegistryMockRecorder) Load(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockBackendRegistry)(nil).Load), arg0, arg1)
}

// Store mocks base method
func (m *MockBackendRegistry) Store(arg0 context.Context, arg1 string, arg2 http.RoundTripper) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Store", arg0, arg1, arg2)
}

// Store indicates an expected call of Store
func (mr *MockBackendRegistryMockRecorder) Store(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockBackendRegistry)(nil).Store), arg0, arg1, arg2)
}
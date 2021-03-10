// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/discovery/interface.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	discovery "github.com/isnlan/coral/pkg/discovery"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockServiceDiscover is a mock of ServiceDiscover interface
type MockServiceDiscover struct {
	ctrl     *gomock.Controller
	recorder *MockServiceDiscoverMockRecorder
}

// MockServiceDiscoverMockRecorder is the mock recorder for MockServiceDiscover
type MockServiceDiscoverMockRecorder struct {
	mock *MockServiceDiscover
}

// NewMockServiceDiscover creates a new mock instance
func NewMockServiceDiscover(ctrl *gomock.Controller) *MockServiceDiscover {
	mock := &MockServiceDiscover{ctrl: ctrl}
	mock.recorder = &MockServiceDiscoverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServiceDiscover) EXPECT() *MockServiceDiscoverMockRecorder {
	return m.recorder
}

// RegisterHealthServer mocks base method
func (m *MockServiceDiscover) RegisterHealthServer(s *grpc.Server) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RegisterHealthServer", s)
}

// RegisterHealthServer indicates an expected call of RegisterHealthServer
func (mr *MockServiceDiscoverMockRecorder) RegisterHealthServer(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterHealthServer", reflect.TypeOf((*MockServiceDiscover)(nil).RegisterHealthServer), s)
}

// ServiceRegister mocks base method
func (m *MockServiceDiscover) ServiceRegister(name, address string, port int, tags ...string) (discovery.Deregister, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{name, address, port}
	for _, a := range tags {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ServiceRegister", varargs...)
	ret0, _ := ret[0].(discovery.Deregister)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ServiceRegister indicates an expected call of ServiceRegister
func (mr *MockServiceDiscoverMockRecorder) ServiceRegister(name, address, port interface{}, tags ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name, address, port}, tags...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceRegister", reflect.TypeOf((*MockServiceDiscover)(nil).ServiceRegister), varargs...)
}

// WatchService mocks base method
func (m *MockServiceDiscover) WatchService(ctx context.Context, name, tag string, ch chan<- []*discovery.ServiceInfo) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WatchService", ctx, name, tag, ch)
}

// WatchService indicates an expected call of WatchService
func (mr *MockServiceDiscoverMockRecorder) WatchService(ctx, name, tag, ch interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchService", reflect.TypeOf((*MockServiceDiscover)(nil).WatchService), ctx, name, tag, ch)
}

// Code generated by MockGen. DO NOT EDIT.
// Source: proto/profile.pb.go

// Package mock_proto is a generated GoMock package.
package mock

import (
	proto "2021_2_LostPointer/internal/microservices/profile/proto"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockProfileClient is a mock of ProfileClient interface.
type MockProfileClient struct {
	ctrl     *gomock.Controller
	recorder *MockProfileClientMockRecorder
}

// MockProfileClientMockRecorder is the mock recorder for MockProfileClient.
type MockProfileClientMockRecorder struct {
	mock *MockProfileClient
}

// NewMockProfileClient creates a new mock instance.
func NewMockProfileClient(ctrl *gomock.Controller) *MockProfileClient {
	mock := &MockProfileClient{ctrl: ctrl}
	mock.recorder = &MockProfileClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfileClient) EXPECT() *MockProfileClientMockRecorder {
	return m.recorder
}

// GetSettings mocks base method.
func (m *MockProfileClient) GetSettings(ctx context.Context, in *proto.GetSettingsOptions, opts ...grpc.CallOption) (*proto.UserSettings, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSettings", varargs...)
	ret0, _ := ret[0].(*proto.UserSettings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSettings indicates an expected call of GetSettings.
func (mr *MockProfileClientMockRecorder) GetSettings(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSettings", reflect.TypeOf((*MockProfileClient)(nil).GetSettings), varargs...)
}

// UpdateSettings mocks base method.
func (m *MockProfileClient) UpdateSettings(ctx context.Context, in *proto.UpdateSettingsOptions, opts ...grpc.CallOption) (*proto.EmptyProfile, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdateSettings", varargs...)
	ret0, _ := ret[0].(*proto.EmptyProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSettings indicates an expected call of UpdateSettings.
func (mr *MockProfileClientMockRecorder) UpdateSettings(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSettings", reflect.TypeOf((*MockProfileClient)(nil).UpdateSettings), varargs...)
}

// MockProfileServer is a mock of ProfileServer interface.
type MockProfileServer struct {
	ctrl     *gomock.Controller
	recorder *MockProfileServerMockRecorder
}

// MockProfileServerMockRecorder is the mock recorder for MockProfileServer.
type MockProfileServerMockRecorder struct {
	mock *MockProfileServer
}

// NewMockProfileServer creates a new mock instance.
func NewMockProfileServer(ctrl *gomock.Controller) *MockProfileServer {
	mock := &MockProfileServer{ctrl: ctrl}
	mock.recorder = &MockProfileServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProfileServer) EXPECT() *MockProfileServerMockRecorder {
	return m.recorder
}

// GetSettings mocks base method.
func (m *MockProfileServer) GetSettings(arg0 context.Context, arg1 *proto.GetSettingsOptions) (*proto.UserSettings, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSettings", arg0, arg1)
	ret0, _ := ret[0].(*proto.UserSettings)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSettings indicates an expected call of GetSettings.
func (mr *MockProfileServerMockRecorder) GetSettings(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSettings", reflect.TypeOf((*MockProfileServer)(nil).GetSettings), arg0, arg1)
}

// UpdateSettings mocks base method.
func (m *MockProfileServer) UpdateSettings(arg0 context.Context, arg1 *proto.UpdateSettingsOptions) (*proto.EmptyProfile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSettings", arg0, arg1)
	ret0, _ := ret[0].(*proto.EmptyProfile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSettings indicates an expected call of UpdateSettings.
func (mr *MockProfileServerMockRecorder) UpdateSettings(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSettings", reflect.TypeOf((*MockProfileServer)(nil).UpdateSettings), arg0, arg1)
}

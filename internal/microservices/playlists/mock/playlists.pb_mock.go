// Code generated by MockGen. DO NOT EDIT.
// Source: proto/playlists.pb.go

// Package mock_proto is a generated GoMock package.
package mock

import (
	proto "2021_2_LostPointer/internal/microservices/playlists/proto"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockPlaylistsClient is a mock of PlaylistsClient interface.
type MockPlaylistsClient struct {
	ctrl     *gomock.Controller
	recorder *MockPlaylistsClientMockRecorder
}

// MockPlaylistsClientMockRecorder is the mock recorder for MockPlaylistsClient.
type MockPlaylistsClientMockRecorder struct {
	mock *MockPlaylistsClient
}

// NewMockPlaylistsClient creates a new mock instance.
func NewMockPlaylistsClient(ctrl *gomock.Controller) *MockPlaylistsClient {
	mock := &MockPlaylistsClient{ctrl: ctrl}
	mock.recorder = &MockPlaylistsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlaylistsClient) EXPECT() *MockPlaylistsClientMockRecorder {
	return m.recorder
}

// AddTrack mocks base method.
func (m *MockPlaylistsClient) AddTrack(ctx context.Context, in *proto.AddTrackOptions, opts ...grpc.CallOption) (*proto.AddTrackResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddTrack", varargs...)
	ret0, _ := ret[0].(*proto.AddTrackResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTrack indicates an expected call of AddTrack.
func (mr *MockPlaylistsClientMockRecorder) AddTrack(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrack", reflect.TypeOf((*MockPlaylistsClient)(nil).AddTrack), varargs...)
}

// CreatePlaylist mocks base method.
func (m *MockPlaylistsClient) CreatePlaylist(ctx context.Context, in *proto.CreatePlaylistOptions, opts ...grpc.CallOption) (*proto.CreatePlaylistResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CreatePlaylist", varargs...)
	ret0, _ := ret[0].(*proto.CreatePlaylistResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePlaylist indicates an expected call of CreatePlaylist.
func (mr *MockPlaylistsClientMockRecorder) CreatePlaylist(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePlaylist", reflect.TypeOf((*MockPlaylistsClient)(nil).CreatePlaylist), varargs...)
}

// DeletePlaylist mocks base method.
func (m *MockPlaylistsClient) DeletePlaylist(ctx context.Context, in *proto.DeletePlaylistOptions, opts ...grpc.CallOption) (*proto.DeletePlaylistResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeletePlaylist", varargs...)
	ret0, _ := ret[0].(*proto.DeletePlaylistResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletePlaylist indicates an expected call of DeletePlaylist.
func (mr *MockPlaylistsClientMockRecorder) DeletePlaylist(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlaylist", reflect.TypeOf((*MockPlaylistsClient)(nil).DeletePlaylist), varargs...)
}

// DeletePlaylistArtwork mocks base method.
func (m *MockPlaylistsClient) DeletePlaylistArtwork(ctx context.Context, in *proto.DeletePlaylistArtworkOptions, opts ...grpc.CallOption) (*proto.DeletePlaylistArtworkResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeletePlaylistArtwork", varargs...)
	ret0, _ := ret[0].(*proto.DeletePlaylistArtworkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletePlaylistArtwork indicates an expected call of DeletePlaylistArtwork.
func (mr *MockPlaylistsClientMockRecorder) DeletePlaylistArtwork(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlaylistArtwork", reflect.TypeOf((*MockPlaylistsClient)(nil).DeletePlaylistArtwork), varargs...)
}

// DeleteTrack mocks base method.
func (m *MockPlaylistsClient) DeleteTrack(ctx context.Context, in *proto.DeleteTrackOptions, opts ...grpc.CallOption) (*proto.DeleteTrackResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteTrack", varargs...)
	ret0, _ := ret[0].(*proto.DeleteTrackResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTrack indicates an expected call of DeleteTrack.
func (mr *MockPlaylistsClientMockRecorder) DeleteTrack(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrack", reflect.TypeOf((*MockPlaylistsClient)(nil).DeleteTrack), varargs...)
}

// UpdatePlaylist mocks base method.
func (m *MockPlaylistsClient) UpdatePlaylist(ctx context.Context, in *proto.UpdatePlaylistOptions, opts ...grpc.CallOption) (*proto.UpdatePlaylistResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpdatePlaylist", varargs...)
	ret0, _ := ret[0].(*proto.UpdatePlaylistResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePlaylist indicates an expected call of UpdatePlaylist.
func (mr *MockPlaylistsClientMockRecorder) UpdatePlaylist(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePlaylist", reflect.TypeOf((*MockPlaylistsClient)(nil).UpdatePlaylist), varargs...)
}

// MockPlaylistsServer is a mock of PlaylistsServer interface.
type MockPlaylistsServer struct {
	ctrl     *gomock.Controller
	recorder *MockPlaylistsServerMockRecorder
}

// MockPlaylistsServerMockRecorder is the mock recorder for MockPlaylistsServer.
type MockPlaylistsServerMockRecorder struct {
	mock *MockPlaylistsServer
}

// NewMockPlaylistsServer creates a new mock instance.
func NewMockPlaylistsServer(ctrl *gomock.Controller) *MockPlaylistsServer {
	mock := &MockPlaylistsServer{ctrl: ctrl}
	mock.recorder = &MockPlaylistsServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlaylistsServer) EXPECT() *MockPlaylistsServerMockRecorder {
	return m.recorder
}

// AddTrack mocks base method.
func (m *MockPlaylistsServer) AddTrack(arg0 context.Context, arg1 *proto.AddTrackOptions) (*proto.AddTrackResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTrack", arg0, arg1)
	ret0, _ := ret[0].(*proto.AddTrackResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTrack indicates an expected call of AddTrack.
func (mr *MockPlaylistsServerMockRecorder) AddTrack(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrack", reflect.TypeOf((*MockPlaylistsServer)(nil).AddTrack), arg0, arg1)
}

// CreatePlaylist mocks base method.
func (m *MockPlaylistsServer) CreatePlaylist(arg0 context.Context, arg1 *proto.CreatePlaylistOptions) (*proto.CreatePlaylistResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePlaylist", arg0, arg1)
	ret0, _ := ret[0].(*proto.CreatePlaylistResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePlaylist indicates an expected call of CreatePlaylist.
func (mr *MockPlaylistsServerMockRecorder) CreatePlaylist(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePlaylist", reflect.TypeOf((*MockPlaylistsServer)(nil).CreatePlaylist), arg0, arg1)
}

// DeletePlaylist mocks base method.
func (m *MockPlaylistsServer) DeletePlaylist(arg0 context.Context, arg1 *proto.DeletePlaylistOptions) (*proto.DeletePlaylistResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlaylist", arg0, arg1)
	ret0, _ := ret[0].(*proto.DeletePlaylistResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletePlaylist indicates an expected call of DeletePlaylist.
func (mr *MockPlaylistsServerMockRecorder) DeletePlaylist(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlaylist", reflect.TypeOf((*MockPlaylistsServer)(nil).DeletePlaylist), arg0, arg1)
}

// DeletePlaylistArtwork mocks base method.
func (m *MockPlaylistsServer) DeletePlaylistArtwork(arg0 context.Context, arg1 *proto.DeletePlaylistArtworkOptions) (*proto.DeletePlaylistArtworkResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePlaylistArtwork", arg0, arg1)
	ret0, _ := ret[0].(*proto.DeletePlaylistArtworkResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeletePlaylistArtwork indicates an expected call of DeletePlaylistArtwork.
func (mr *MockPlaylistsServerMockRecorder) DeletePlaylistArtwork(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePlaylistArtwork", reflect.TypeOf((*MockPlaylistsServer)(nil).DeletePlaylistArtwork), arg0, arg1)
}

// DeleteTrack mocks base method.
func (m *MockPlaylistsServer) DeleteTrack(arg0 context.Context, arg1 *proto.DeleteTrackOptions) (*proto.DeleteTrackResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrack", arg0, arg1)
	ret0, _ := ret[0].(*proto.DeleteTrackResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTrack indicates an expected call of DeleteTrack.
func (mr *MockPlaylistsServerMockRecorder) DeleteTrack(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrack", reflect.TypeOf((*MockPlaylistsServer)(nil).DeleteTrack), arg0, arg1)
}

// UpdatePlaylist mocks base method.
func (m *MockPlaylistsServer) UpdatePlaylist(arg0 context.Context, arg1 *proto.UpdatePlaylistOptions) (*proto.UpdatePlaylistResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePlaylist", arg0, arg1)
	ret0, _ := ret[0].(*proto.UpdatePlaylistResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePlaylist indicates an expected call of UpdatePlaylist.
func (mr *MockPlaylistsServerMockRecorder) UpdatePlaylist(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePlaylist", reflect.TypeOf((*MockPlaylistsServer)(nil).UpdatePlaylist), arg0, arg1)
}

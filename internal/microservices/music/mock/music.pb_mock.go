// Code generated by MockGen. DO NOT EDIT.
// Source: proto/music.pb.go

// Package mock_proto is a generated GoMock package.
package mock

import (
	proto "2021_2_LostPointer/internal/microservices/music/proto"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockMusicClient is a mock of MusicClient interface.
type MockMusicClient struct {
	ctrl     *gomock.Controller
	recorder *MockMusicClientMockRecorder
}

// MockMusicClientMockRecorder is the mock recorder for MockMusicClient.
type MockMusicClientMockRecorder struct {
	mock *MockMusicClient
}

// NewMockMusicClient creates a new mock instance.
func NewMockMusicClient(ctrl *gomock.Controller) *MockMusicClient {
	mock := &MockMusicClient{ctrl: ctrl}
	mock.recorder = &MockMusicClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMusicClient) EXPECT() *MockMusicClientMockRecorder {
	return m.recorder
}

// AddTrackToFavorites mocks base method.
func (m *MockMusicClient) AddTrackToFavorites(ctx context.Context, in *proto.AddTrackToFavoritesOptions, opts ...grpc.CallOption) (*proto.AddTrackToFavoritesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddTrackToFavorites", varargs...)
	ret0, _ := ret[0].(*proto.AddTrackToFavoritesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTrackToFavorites indicates an expected call of AddTrackToFavorites.
func (mr *MockMusicClientMockRecorder) AddTrackToFavorites(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrackToFavorites", reflect.TypeOf((*MockMusicClient)(nil).AddTrackToFavorites), varargs...)
}

// AlbumPage mocks base method.
func (m *MockMusicClient) AlbumPage(ctx context.Context, in *proto.AlbumPageOptions, opts ...grpc.CallOption) (*proto.AlbumPageResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AlbumPage", varargs...)
	ret0, _ := ret[0].(*proto.AlbumPageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AlbumPage indicates an expected call of AlbumPage.
func (mr *MockMusicClientMockRecorder) AlbumPage(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AlbumPage", reflect.TypeOf((*MockMusicClient)(nil).AlbumPage), varargs...)
}

// ArtistProfile mocks base method.
func (m *MockMusicClient) ArtistProfile(ctx context.Context, in *proto.ArtistProfileOptions, opts ...grpc.CallOption) (*proto.Artist, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ArtistProfile", varargs...)
	ret0, _ := ret[0].(*proto.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ArtistProfile indicates an expected call of ArtistProfile.
func (mr *MockMusicClientMockRecorder) ArtistProfile(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArtistProfile", reflect.TypeOf((*MockMusicClient)(nil).ArtistProfile), varargs...)
}

// DeleteTrackFromFavorites mocks base method.
func (m *MockMusicClient) DeleteTrackFromFavorites(ctx context.Context, in *proto.DeleteTrackFromFavoritesOptions, opts ...grpc.CallOption) (*proto.DeleteTrackFromFavoritesResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteTrackFromFavorites", varargs...)
	ret0, _ := ret[0].(*proto.DeleteTrackFromFavoritesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTrackFromFavorites indicates an expected call of DeleteTrackFromFavorites.
func (mr *MockMusicClientMockRecorder) DeleteTrackFromFavorites(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrackFromFavorites", reflect.TypeOf((*MockMusicClient)(nil).DeleteTrackFromFavorites), varargs...)
}

// Find mocks base method.
func (m *MockMusicClient) Find(ctx context.Context, in *proto.FindOptions, opts ...grpc.CallOption) (*proto.FindResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Find", varargs...)
	ret0, _ := ret[0].(*proto.FindResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockMusicClientMockRecorder) Find(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockMusicClient)(nil).Find), varargs...)
}

// GetFavoriteTracks mocks base method.
func (m *MockMusicClient) GetFavoriteTracks(ctx context.Context, in *proto.UserFavoritesOptions, opts ...grpc.CallOption) (*proto.Tracks, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFavoriteTracks", varargs...)
	ret0, _ := ret[0].(*proto.Tracks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavoriteTracks indicates an expected call of GetFavoriteTracks.
func (mr *MockMusicClientMockRecorder) GetFavoriteTracks(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavoriteTracks", reflect.TypeOf((*MockMusicClient)(nil).GetFavoriteTracks), varargs...)
}

// IncrementListenCount mocks base method.
func (m *MockMusicClient) IncrementListenCount(ctx context.Context, in *proto.IncrementListenCountOptions, opts ...grpc.CallOption) (*proto.IncrementListenCountEmpty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IncrementListenCount", varargs...)
	ret0, _ := ret[0].(*proto.IncrementListenCountEmpty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IncrementListenCount indicates an expected call of IncrementListenCount.
func (mr *MockMusicClientMockRecorder) IncrementListenCount(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementListenCount", reflect.TypeOf((*MockMusicClient)(nil).IncrementListenCount), varargs...)
}

// PlaylistPage mocks base method.
func (m *MockMusicClient) PlaylistPage(ctx context.Context, in *proto.PlaylistPageOptions, opts ...grpc.CallOption) (*proto.PlaylistPageResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PlaylistPage", varargs...)
	ret0, _ := ret[0].(*proto.PlaylistPageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlaylistPage indicates an expected call of PlaylistPage.
func (mr *MockMusicClientMockRecorder) PlaylistPage(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlaylistPage", reflect.TypeOf((*MockMusicClient)(nil).PlaylistPage), varargs...)
}

// RandomAlbums mocks base method.
func (m *MockMusicClient) RandomAlbums(ctx context.Context, in *proto.RandomAlbumsOptions, opts ...grpc.CallOption) (*proto.Albums, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RandomAlbums", varargs...)
	ret0, _ := ret[0].(*proto.Albums)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RandomAlbums indicates an expected call of RandomAlbums.
func (mr *MockMusicClientMockRecorder) RandomAlbums(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomAlbums", reflect.TypeOf((*MockMusicClient)(nil).RandomAlbums), varargs...)
}

// RandomArtists mocks base method.
func (m *MockMusicClient) RandomArtists(ctx context.Context, in *proto.RandomArtistsOptions, opts ...grpc.CallOption) (*proto.Artists, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RandomArtists", varargs...)
	ret0, _ := ret[0].(*proto.Artists)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RandomArtists indicates an expected call of RandomArtists.
func (mr *MockMusicClientMockRecorder) RandomArtists(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomArtists", reflect.TypeOf((*MockMusicClient)(nil).RandomArtists), varargs...)
}

// RandomTracks mocks base method.
func (m *MockMusicClient) RandomTracks(ctx context.Context, in *proto.RandomTracksOptions, opts ...grpc.CallOption) (*proto.Tracks, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "RandomTracks", varargs...)
	ret0, _ := ret[0].(*proto.Tracks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RandomTracks indicates an expected call of RandomTracks.
func (mr *MockMusicClientMockRecorder) RandomTracks(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomTracks", reflect.TypeOf((*MockMusicClient)(nil).RandomTracks), varargs...)
}

// UserPlaylists mocks base method.
func (m *MockMusicClient) UserPlaylists(ctx context.Context, in *proto.UserPlaylistsOptions, opts ...grpc.CallOption) (*proto.PlaylistsData, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UserPlaylists", varargs...)
	ret0, _ := ret[0].(*proto.PlaylistsData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserPlaylists indicates an expected call of UserPlaylists.
func (mr *MockMusicClientMockRecorder) UserPlaylists(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserPlaylists", reflect.TypeOf((*MockMusicClient)(nil).UserPlaylists), varargs...)
}

// MockMusicServer is a mock of MusicServer interface.
type MockMusicServer struct {
	ctrl     *gomock.Controller
	recorder *MockMusicServerMockRecorder
}

// MockMusicServerMockRecorder is the mock recorder for MockMusicServer.
type MockMusicServerMockRecorder struct {
	mock *MockMusicServer
}

// NewMockMusicServer creates a new mock instance.
func NewMockMusicServer(ctrl *gomock.Controller) *MockMusicServer {
	mock := &MockMusicServer{ctrl: ctrl}
	mock.recorder = &MockMusicServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMusicServer) EXPECT() *MockMusicServerMockRecorder {
	return m.recorder
}

// AddTrackToFavorites mocks base method.
func (m *MockMusicServer) AddTrackToFavorites(arg0 context.Context, arg1 *proto.AddTrackToFavoritesOptions) (*proto.AddTrackToFavoritesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddTrackToFavorites", arg0, arg1)
	ret0, _ := ret[0].(*proto.AddTrackToFavoritesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddTrackToFavorites indicates an expected call of AddTrackToFavorites.
func (mr *MockMusicServerMockRecorder) AddTrackToFavorites(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddTrackToFavorites", reflect.TypeOf((*MockMusicServer)(nil).AddTrackToFavorites), arg0, arg1)
}

// AlbumPage mocks base method.
func (m *MockMusicServer) AlbumPage(arg0 context.Context, arg1 *proto.AlbumPageOptions) (*proto.AlbumPageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AlbumPage", arg0, arg1)
	ret0, _ := ret[0].(*proto.AlbumPageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AlbumPage indicates an expected call of AlbumPage.
func (mr *MockMusicServerMockRecorder) AlbumPage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AlbumPage", reflect.TypeOf((*MockMusicServer)(nil).AlbumPage), arg0, arg1)
}

// ArtistProfile mocks base method.
func (m *MockMusicServer) ArtistProfile(arg0 context.Context, arg1 *proto.ArtistProfileOptions) (*proto.Artist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ArtistProfile", arg0, arg1)
	ret0, _ := ret[0].(*proto.Artist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ArtistProfile indicates an expected call of ArtistProfile.
func (mr *MockMusicServerMockRecorder) ArtistProfile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ArtistProfile", reflect.TypeOf((*MockMusicServer)(nil).ArtistProfile), arg0, arg1)
}

// DeleteTrackFromFavorites mocks base method.
func (m *MockMusicServer) DeleteTrackFromFavorites(arg0 context.Context, arg1 *proto.DeleteTrackFromFavoritesOptions) (*proto.DeleteTrackFromFavoritesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrackFromFavorites", arg0, arg1)
	ret0, _ := ret[0].(*proto.DeleteTrackFromFavoritesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteTrackFromFavorites indicates an expected call of DeleteTrackFromFavorites.
func (mr *MockMusicServerMockRecorder) DeleteTrackFromFavorites(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrackFromFavorites", reflect.TypeOf((*MockMusicServer)(nil).DeleteTrackFromFavorites), arg0, arg1)
}

// Find mocks base method.
func (m *MockMusicServer) Find(arg0 context.Context, arg1 *proto.FindOptions) (*proto.FindResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", arg0, arg1)
	ret0, _ := ret[0].(*proto.FindResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockMusicServerMockRecorder) Find(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockMusicServer)(nil).Find), arg0, arg1)
}

// GetFavoriteTracks mocks base method.
func (m *MockMusicServer) GetFavoriteTracks(arg0 context.Context, arg1 *proto.UserFavoritesOptions) (*proto.Tracks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFavoriteTracks", arg0, arg1)
	ret0, _ := ret[0].(*proto.Tracks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFavoriteTracks indicates an expected call of GetFavoriteTracks.
func (mr *MockMusicServerMockRecorder) GetFavoriteTracks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFavoriteTracks", reflect.TypeOf((*MockMusicServer)(nil).GetFavoriteTracks), arg0, arg1)
}

// IncrementListenCount mocks base method.
func (m *MockMusicServer) IncrementListenCount(arg0 context.Context, arg1 *proto.IncrementListenCountOptions) (*proto.IncrementListenCountEmpty, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementListenCount", arg0, arg1)
	ret0, _ := ret[0].(*proto.IncrementListenCountEmpty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IncrementListenCount indicates an expected call of IncrementListenCount.
func (mr *MockMusicServerMockRecorder) IncrementListenCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementListenCount", reflect.TypeOf((*MockMusicServer)(nil).IncrementListenCount), arg0, arg1)
}

// PlaylistPage mocks base method.
func (m *MockMusicServer) PlaylistPage(arg0 context.Context, arg1 *proto.PlaylistPageOptions) (*proto.PlaylistPageResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlaylistPage", arg0, arg1)
	ret0, _ := ret[0].(*proto.PlaylistPageResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlaylistPage indicates an expected call of PlaylistPage.
func (mr *MockMusicServerMockRecorder) PlaylistPage(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlaylistPage", reflect.TypeOf((*MockMusicServer)(nil).PlaylistPage), arg0, arg1)
}

// RandomAlbums mocks base method.
func (m *MockMusicServer) RandomAlbums(arg0 context.Context, arg1 *proto.RandomAlbumsOptions) (*proto.Albums, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RandomAlbums", arg0, arg1)
	ret0, _ := ret[0].(*proto.Albums)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RandomAlbums indicates an expected call of RandomAlbums.
func (mr *MockMusicServerMockRecorder) RandomAlbums(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomAlbums", reflect.TypeOf((*MockMusicServer)(nil).RandomAlbums), arg0, arg1)
}

// RandomArtists mocks base method.
func (m *MockMusicServer) RandomArtists(arg0 context.Context, arg1 *proto.RandomArtistsOptions) (*proto.Artists, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RandomArtists", arg0, arg1)
	ret0, _ := ret[0].(*proto.Artists)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RandomArtists indicates an expected call of RandomArtists.
func (mr *MockMusicServerMockRecorder) RandomArtists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomArtists", reflect.TypeOf((*MockMusicServer)(nil).RandomArtists), arg0, arg1)
}

// RandomTracks mocks base method.
func (m *MockMusicServer) RandomTracks(arg0 context.Context, arg1 *proto.RandomTracksOptions) (*proto.Tracks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RandomTracks", arg0, arg1)
	ret0, _ := ret[0].(*proto.Tracks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RandomTracks indicates an expected call of RandomTracks.
func (mr *MockMusicServerMockRecorder) RandomTracks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RandomTracks", reflect.TypeOf((*MockMusicServer)(nil).RandomTracks), arg0, arg1)
}

// UserPlaylists mocks base method.
func (m *MockMusicServer) UserPlaylists(arg0 context.Context, arg1 *proto.UserPlaylistsOptions) (*proto.PlaylistsData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserPlaylists", arg0, arg1)
	ret0, _ := ret[0].(*proto.PlaylistsData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserPlaylists indicates an expected call of UserPlaylists.
func (mr *MockMusicServerMockRecorder) UserPlaylists(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserPlaylists", reflect.TypeOf((*MockMusicServer)(nil).UserPlaylists), arg0, arg1)
}

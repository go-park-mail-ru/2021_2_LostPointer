package usecase

/*
import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/music/mock"
	"2021_2_LostPointer/internal/microservices/music/proto"
)

func TestMusicService_RandomTracks(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.RandomTracksOptions
		expected    *proto.Tracks
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				RandomTracksFunc: func(int64, bool) (*proto.Tracks, error) {
					return &proto.Tracks{
						Tracks: []*proto.Track{},
					}, nil
				},
			},
			input: &proto.RandomTracksOptions{
				Amount:       1,
				IsAuthorized: true,
			},
			expected: &proto.Tracks{
				Tracks: []*proto.Track{},
			},
		},
		{
			name: "Error 500. mock.RandomTracks returned error",
			storageMock: &mock.MockStorage{
				RandomTracksFunc: func(int64, bool) (*proto.Tracks, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.RandomTracksOptions{
				Amount:       1,
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.RandomTracks(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestMusicService_RandomAlbums(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.RandomAlbumsOptions
		expected    *proto.Albums
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				RandomAlbumsFunc: func(int64) (*proto.Albums, error) {
					return &proto.Albums{
						Albums: []*proto.Album{},
					}, nil
				},
			},
			input: &proto.RandomAlbumsOptions{Amount: 5},
			expected: &proto.Albums{
				Albums: []*proto.Album{},
			},
		},
		{
			name: "Error 500. mock.RandomAlbums returned error",
			storageMock: &mock.MockStorage{
				RandomAlbumsFunc: func(int64) (*proto.Albums, error) {
					return nil, errors.New("error")
				},
			},
			input:       &proto.RandomAlbumsOptions{Amount: 5},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.RandomAlbums(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestMusicService_RandomArtists(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.RandomArtistsOptions
		expected    *proto.Artists
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				RandomArtistsFunc: func(int64) (*proto.Artists, error) {
					return &proto.Artists{
						Artists: []*proto.Artist{},
					}, nil
				},
			},
			input: &proto.RandomArtistsOptions{Amount: 5},
			expected: &proto.Artists{
				Artists: []*proto.Artist{},
			},
		},
		{
			name: "Error 500. mock.RandomArtists returned error",
			storageMock: &mock.MockStorage{
				RandomArtistsFunc: func(int64) (*proto.Artists, error) {
					return nil, errors.New("error")
				},
			},
			input:       &proto.RandomArtistsOptions{Amount: 5},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.RandomArtists(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestMusicService_ArtistProfile(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.ArtistProfileOptions
		expected    *proto.Artist
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				ArtistInfoFunc: func(int64) (*proto.Artist, error) {
					return &proto.Artist{}, nil
				},
				ArtistTracksFunc: func(int64, bool, int64) ([]*proto.Track, error) {
					return nil, nil
				},
				ArtistAlbumsFunc: func(int64, int64) ([]*proto.Album, error) {
					return nil, nil
				},
			},
			input: &proto.ArtistProfileOptions{
				ArtistID:     1,
				IsAuthorized: true,
			},
			expected: &proto.Artist{},
		},
		{
			name: "Error 500. mock.ArtistInfo returned error",
			storageMock: &mock.MockStorage{
				ArtistInfoFunc: func(int64) (*proto.Artist, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.ArtistProfileOptions{
				ArtistID:     1,
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.ArtistTracks returned error",
			storageMock: &mock.MockStorage{
				ArtistInfoFunc: func(int64) (*proto.Artist, error) {
					return &proto.Artist{}, nil
				},
				ArtistTracksFunc: func(int64, bool, int64) ([]*proto.Track, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.ArtistProfileOptions{
				ArtistID:     1,
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.ArtistAlbums returned error",
			storageMock: &mock.MockStorage{
				ArtistInfoFunc: func(int64) (*proto.Artist, error) {
					return &proto.Artist{}, nil
				},
				ArtistTracksFunc: func(int64, bool, int64) ([]*proto.Track, error) {
					return nil, nil
				},
				ArtistAlbumsFunc: func(int64, int64) ([]*proto.Album, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.ArtistProfileOptions{
				ArtistID:     1,
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.ArtistProfile(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestMusicService_IncrementListenCount(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.IncrementListenCountOptions
		expected    *proto.IncrementListenCountEmpty
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				IncrementListenCountFunc: func(int64) error {
					return nil
				},
			},
			input:    &proto.IncrementListenCountOptions{ID: 1},
			expected: &proto.IncrementListenCountEmpty{},
		},
		{
			name: "Error 500. mock.IncrementListenCount returned error",
			storageMock: &mock.MockStorage{
				IncrementListenCountFunc: func(int64) error {
					return errors.New("error")
				},
			},
			input:       &proto.IncrementListenCountOptions{ID: 1},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.IncrementListenCount(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestMusicService_AlbumPage(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.AlbumPageOptions
		expected    *proto.AlbumPageResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				AlbumDataFunc: func(int64) (*proto.AlbumPageResponse, error) {
					return &proto.AlbumPageResponse{}, nil
				},
				AlbumTracksFunc: func(int64, bool) ([]*proto.AlbumTrack, error) {
					return nil, nil
				},
			},
			input: &proto.AlbumPageOptions{
				AlbumID:      1,
				IsAuthorized: true,
			},
			expected: &proto.AlbumPageResponse{},
		},
		{
			name: "Error 500. mock.AlbumData returned error",
			storageMock: &mock.MockStorage{
				AlbumDataFunc: func(int64) (*proto.AlbumPageResponse, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.AlbumPageOptions{
				AlbumID:      1,
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.AlbumTracks returned error",
			storageMock: &mock.MockStorage{
				AlbumDataFunc: func(int64) (*proto.AlbumPageResponse, error) {
					return &proto.AlbumPageResponse{}, nil
				},
				AlbumTracksFunc: func(int64, bool) ([]*proto.AlbumTrack, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.AlbumPageOptions{
				AlbumID:      1,
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.AlbumPage(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestMusicService_UserPlaylists(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.UserPlaylistsOptions
		expected    *proto.PlaylistsData
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				UserPlaylistsFunc: func(int64) ([]*proto.PlaylistData, error) {
					return []*proto.PlaylistData{}, nil

				},
			},
			input: &proto.UserPlaylistsOptions{UserID: 1},
			expected: &proto.PlaylistsData{
				Playlists: []*proto.PlaylistData{},
			},
		},
		{
			name: "Error 500. mock.UserPlaylists returned error",
			storageMock: &mock.MockStorage{
				UserPlaylistsFunc: func(int64) ([]*proto.PlaylistData, error) {
					return nil, errors.New("error")
				},
			},
			input:       &proto.UserPlaylistsOptions{UserID: 1},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.UserPlaylists(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestMusicService_Find(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.FindOptions
		expected    *proto.FindResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Success, contain returned true",
			storageMock: &mock.MockStorage{
				FindTracksByFullWordFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{
						{
							ID:          0,
							Title:       "",
							Explicit:    false,
							Genre:       "",
							Number:      0,
							File:        "",
							ListenCount: 0,
							Duration:    0,
							Lossless:    false,
							Album:       nil,
							Artist:      nil,
						},
					}, nil
				},
				FindTracksByPartialFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{
						{
							ID:          0,
							Title:       "",
							Explicit:    false,
							Genre:       "",
							Number:      0,
							File:        "",
							ListenCount: 0,
							Duration:    0,
							Lossless:    false,
							Album:       nil,
							Artist:      nil,
						},
					}, nil
				},
				FindArtistsFunc: func(string) ([]*proto.Artist, error) {
					return []*proto.Artist{}, nil
				},
				FindAlbumsFunc: func(string) ([]*proto.Album, error) {
					return []*proto.Album{}, nil
				},
			},
			input: &proto.FindOptions{
				Text:         "lahaine",
				IsAuthorized: true,
			},
			expected: &proto.FindResponse{
				Tracks: []*proto.Track{
					{
						ID:          0,
						Title:       "",
						Explicit:    false,
						Genre:       "",
						Number:      0,
						File:        "",
						ListenCount: 0,
						Duration:    0,
						Lossless:    false,
						Album:       nil,
						Artist:      nil,
					},
				},
				Albums:  []*proto.Album{},
				Artists: []*proto.Artist{},
			},
		},
		{
			name: "Success, contain returned false",
			storageMock: &mock.MockStorage{
				FindTracksByFullWordFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{}, nil
				},
				FindTracksByPartialFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{
						{
							ID:          0,
							Title:       "",
							Explicit:    false,
							Genre:       "",
							Number:      0,
							File:        "",
							ListenCount: 0,
							Duration:    0,
							Lossless:    false,
							Album:       nil,
							Artist:      nil,
						},
					}, nil
				},
				FindArtistsFunc: func(string) ([]*proto.Artist, error) {
					return []*proto.Artist{}, nil
				},
				FindAlbumsFunc: func(string) ([]*proto.Album, error) {
					return []*proto.Album{}, nil
				},
			},
			input: &proto.FindOptions{
				Text:         "lahaine",
				IsAuthorized: true,
			},
			expected: &proto.FindResponse{
				Tracks: []*proto.Track{
					{
						ID:          0,
						Title:       "",
						Explicit:    false,
						Genre:       "",
						Number:      0,
						File:        "",
						ListenCount: 0,
						Duration:    0,
						Lossless:    false,
						Album:       nil,
						Artist:      nil,
					},
				},
				Albums:  []*proto.Album{},
				Artists: []*proto.Artist{},
			},
		},
		{
			name: "Error 500. mock.FindTracksByFullWord returned error",
			storageMock: &mock.MockStorage{
				FindTracksByFullWordFunc: func(string, bool) ([]*proto.Track, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.FindOptions{
				Text:         "lahaine",
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.FindTracksByPartial returned error",
			storageMock: &mock.MockStorage{
				FindTracksByFullWordFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{}, nil
				},
				FindTracksByPartialFunc: func(string, bool) ([]*proto.Track, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.FindOptions{
				Text:         "lahaine",
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.FindArtists returned error",
			storageMock: &mock.MockStorage{
				FindTracksByFullWordFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{}, nil
				},
				FindTracksByPartialFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{}, nil
				},
				FindArtistsFunc: func(string) ([]*proto.Artist, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.FindOptions{
				Text:         "lahaine",
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.FindAlbums returned error",
			storageMock: &mock.MockStorage{
				FindTracksByFullWordFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{}, nil
				},
				FindTracksByPartialFunc: func(string, bool) ([]*proto.Track, error) {
					return []*proto.Track{}, nil
				},
				FindArtistsFunc: func(string) ([]*proto.Artist, error) {
					return []*proto.Artist{}, nil
				},
				FindAlbumsFunc: func(string) ([]*proto.Album, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.FindOptions{
				Text:         "lahaine",
				IsAuthorized: true,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.Find(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}

func TestMusicService_PlaylistPage(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.PlaylistPageOptions
		expected    *proto.PlaylistPageResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsPlaylistOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				PlaylistInfoFunc: func(int64) (*proto.PlaylistData, error) {
					return &proto.PlaylistData{}, nil
				},
				PlaylistTracksFunc: func(int64) ([]*proto.Track, error) {
					return []*proto.Track{}, nil
				},
				IsPlaylistPublicFunc: func(int64) (bool, error) {
					return true, nil
				},
			},
			input: &proto.PlaylistPageOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expected: &proto.PlaylistPageResponse{
				Tracks: []*proto.Track{},
				IsOwn:  true,
			},
		},
		{
			name: "Error 500. mock.DoesPlaylistExist returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.PlaylistPageOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 404. mock.DoesPlaylistExist returned false, playlist does not exist",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, nil
				},
			},
			input: &proto.PlaylistPageOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.NotFound, constants.PlaylistNotFoundMessage),
		},
		{
			name: "Error 500. mock.IsPlaylistOwner returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsPlaylistOwnerFunc: func(int64, int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.PlaylistPageOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 403. mock.IsPlaylistOwner returned false, user is not playlist owner",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsPlaylistOwnerFunc: func(int64, int64) (bool, error) {
					return false, nil
				},
				IsPlaylistPublicFunc: func(int64) (bool, error) {
					return false, nil
				},
				PlaylistInfoFunc: func(int64) (*proto.PlaylistData, error) {
					return &proto.PlaylistData{}, nil
				},
				PlaylistTracksFunc: func(int64) ([]*proto.Track, error) {
					return []*proto.Track{}, nil
				},
			},
			input: &proto.PlaylistPageOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage),
		},
		{
			name: "Error 500. mock.PlaylistInfo returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsPlaylistOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				PlaylistInfoFunc: func(int64) (*proto.PlaylistData, error) {
					return nil, errors.New("error")
				},
				IsPlaylistPublicFunc: func(int64) (bool, error) {
					return true, nil
				},
			},
			input: &proto.PlaylistPageOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.PlaylistTracks returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsPlaylistOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				PlaylistInfoFunc: func(int64) (*proto.PlaylistData, error) {
					return &proto.PlaylistData{}, nil
				},
				PlaylistTracksFunc: func(int64) ([]*proto.Track, error) {
					return nil, errors.New("error")
				},
				IsPlaylistPublicFunc: func(int64) (bool, error) {
					return true, nil
				},
			},
			input: &proto.PlaylistPageOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewMusicService(currentTest.storageMock)

			res, err := storage.PlaylistPage(context.Background(), currentTest.input)
			if currentTest.expectedErr {
				assert.Error(t, err)
				assert.Equal(t, err, currentTest.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, currentTest.expected, res)
			}
		})
	}
}
*/

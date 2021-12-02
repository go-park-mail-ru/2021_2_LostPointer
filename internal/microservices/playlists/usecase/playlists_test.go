package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/microservices/playlists/mock"
	"2021_2_LostPointer/internal/microservices/playlists/proto"
)

const testImageFilename = "lahaine.webp"

func TestPlaylistsService_CreatePlaylist(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.CreatePlaylistOptions
		expected    *proto.CreatePlaylistResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				CreatePlaylistFunc: func(int64, string, string, string, bool) (*proto.CreatePlaylistResponse, error) {
					return &proto.CreatePlaylistResponse{}, nil
				},
			},
			input: &proto.CreatePlaylistOptions{
				UserID:       1,
				Title:        "LaHaine A State Of Trance",
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expected: &proto.CreatePlaylistResponse{},
		},
		{
			name:        "Error 400. Title is not valid",
			storageMock: &mock.MockStorage{},
			input: &proto.CreatePlaylistOptions{
				UserID:       1,
				Title:        "A",
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, "The length of title must be from 3 to 30 characters"),
		},
		{
			name: "Error 500. Title is not valid",
			storageMock: &mock.MockStorage{
				CreatePlaylistFunc: func(int64, string, string, string, bool) (*proto.CreatePlaylistResponse, error) {
					return nil, errors.New("error")
				},
			},
			input: &proto.CreatePlaylistOptions{
				UserID:       1,
				Title:        "LaHaine A State Of Trance",
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewPlaylistsService(currentTest.storageMock)

			res, err := storage.CreatePlaylist(context.Background(), currentTest.input)
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

func TestPlaylistsService_DeletePlaylist(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.DeletePlaylistOptions
		expected    *proto.DeletePlaylistResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				DeletePlaylistFunc: func(int64) error {
					return nil
				},
			},
			input: &proto.DeletePlaylistOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expected: &proto.DeletePlaylistResponse{OldArtworkFilename: testImageFilename},
		},
		{
			name: "Error 500. mock.DoesPlaylistExist returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.DeletePlaylistOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. Playlist does not exist",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, nil
				},
			},
			input: &proto.DeletePlaylistOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.NotFound, constants.PlaylistNotFoundMessage),
		},
		{
			name: "Error 500. mock.IsOwner returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.DeletePlaylistOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 403. User is not playlist owner",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, nil
				},
			},
			input: &proto.DeletePlaylistOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage),
		},
		{
			name: "Error 500. mock.GetOldPlaylistSettings returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return "", errors.New("error")
				},
			},
			input: &proto.DeletePlaylistOptions{
				PlaylistID: 1,
				UserID:     1,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.DeletePlaylist returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				DeletePlaylistFunc: func(int64) error {
					return errors.New("error")
				},
			},
			input: &proto.DeletePlaylistOptions{
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
			storage := NewPlaylistsService(currentTest.storageMock)
			res, err := storage.DeletePlaylist(context.Background(), currentTest.input)
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

func TestPlaylistsService_UpdatePlaylist(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.UpdatePlaylistOptions
		expected    *proto.UpdatePlaylistResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Successfully updated title",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				UpdatePlaylistTitleFunc: func(int64, string) error {
					return nil
				},
				UpdatePlaylistAccessFunc: func(int64, bool) error {
					return nil
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "LaHaine new ASOT",
				UserID:       0,
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expected: &proto.UpdatePlaylistResponse{
				OldArtworkFilename: "",
				ArtworkColor:       "",
			},
		},
		{
			name: "Error 500. mock.DoesPlaylistExist returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "LaHaine new ASOT",
				UserID:       0,
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. Playlist does not exist",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, nil
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "LaHaine new ASOT",
				UserID:       0,
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.NotFound, constants.PlaylistNotFoundMessage),
		},
		{
			name: "Error 500. mock.IsOwner returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "LaHaine new ASOT",
				UserID:       0,
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 403. User is not playlist owner",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, nil
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "LaHaine new ASOT",
				UserID:       0,
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage),
		},
		{
			name: "Error 400. New title is not valid",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "La",
				UserID:       0,
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, "The length of title must be from 3 to 30 characters"),
		},
		{
			name: "Error 500. mock.UpdatePlaylistTitle returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				UpdatePlaylistTitleFunc: func(int64, string) error {
					return errors.New("error")
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "LaHaine",
				UserID:       0,
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.UpdatePlaylistAccess returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				UpdatePlaylistAccessFunc: func(int64, bool) error {
					return errors.New("error")
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "",
				UserID:       0,
				Artwork:      "",
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Successfully updated artwork",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				UpdatePlaylistAccessFunc: func(int64, bool) error {
					return nil
				},
				UpdatePlaylistArtworkFunc: func(int64, string, string) error {
					return nil
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "",
				UserID:       0,
				Artwork:      testImageFilename,
				ArtworkColor: "",
				IsPublic:     false,
			},
			expected: &proto.UpdatePlaylistResponse{OldArtworkFilename: testImageFilename},
		},
		{
			name: "Error 500. mock.UpdatePlaylistArtwork returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				UpdatePlaylistAccessFunc: func(int64, bool) error {
					return nil
				},
				UpdatePlaylistArtworkFunc: func(int64, string, string) error {
					return errors.New("error")
				},
			},
			input: &proto.UpdatePlaylistOptions{
				PlaylistID:   1,
				Title:        "",
				UserID:       0,
				Artwork:      testImageFilename,
				ArtworkColor: "",
				IsPublic:     false,
			},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewPlaylistsService(currentTest.storageMock)

			res, err := storage.UpdatePlaylist(context.Background(), currentTest.input)
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

func TestPlaylistsService_AddTrack(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.AddTrackOptions
		expected    *proto.AddTrackResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				IsAddedFunc: func(int64, int64) (bool, error) {
					return false, nil
				},
				AddTrackFunc: func(int64, int64) error {
					return nil
				},
			},
			input:    &proto.AddTrackOptions{},
			expected: &proto.AddTrackResponse{},
		},
		{
			name: "Error 500. mock.DoesPlaylistExist returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input:       &proto.AddTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. Playlist does not exist",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, nil
				},
			},
			input:       &proto.AddTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.NotFound, constants.PlaylistNotFoundMessage),
		},
		{
			name: "Error 500. mock.IsOwner returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input:       &proto.AddTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 403. User is not playlist owner",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, nil
				},
			},
			input:       &proto.AddTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage),
		},
		{
			name: "Error 400. Track is already added",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				IsAddedFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
			},
			input:       &proto.AddTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.InvalidArgument, constants.TrackAlreadyInPlaylistMessage),
		},
		{
			name: "Error 500. mock.IsAdded returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				IsAddedFunc: func(int64, int64) (bool, error) {
					return true, errors.New("error")
				},
			},
			input:       &proto.AddTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.AddTrack returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				IsAddedFunc: func(int64, int64) (bool, error) {
					return false, nil
				},
				AddTrackFunc: func(int64, int64) error {
					return errors.New("error")
				},
			},
			input:       &proto.AddTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewPlaylistsService(currentTest.storageMock)

			res, err := storage.AddTrack(context.Background(), currentTest.input)
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

func TestPlaylistsService_DeleteTrack(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.DeleteTrackOptions
		expected    *proto.DeleteTrackResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				DeleteTrackFunc: func(int64, int64) error {
					return nil
				},
			},
			input:    &proto.DeleteTrackOptions{},
			expected: &proto.DeleteTrackResponse{},
		},
		{
			name: "Error 500. mock.DoesPlaylistExist returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input:       &proto.DeleteTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. Playlist does not exist",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, nil
				},
			},
			input:       &proto.DeleteTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.NotFound, constants.PlaylistNotFoundMessage),
		},
		{
			name: "Error 500. mock.IsOwner returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input:       &proto.DeleteTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 403. User is not playlist owner",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, nil
				},
			},
			input:       &proto.DeleteTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage),
		},
		{
			name: "Error 500. mock.DeleteTrackFunc returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				DeleteTrackFunc: func(int64, int64) error {
					return errors.New("error")
				},
			},
			input:       &proto.DeleteTrackOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewPlaylistsService(currentTest.storageMock)

			res, err := storage.DeleteTrack(context.Background(), currentTest.input)
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

func TestPlaylistsService_DeletePlaylistArtwork(t *testing.T) {
	tests := []struct {
		name        string
		storageMock *mock.MockStorage
		input       *proto.DeletePlaylistArtworkOptions
		expected    *proto.DeletePlaylistArtworkResponse
		expectedErr bool
		err         error
	}{
		{
			name: "Success",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				DeletePlaylistArtworkFunc: func(int64) error {
					return nil
				},
			},
			input:    &proto.DeletePlaylistArtworkOptions{},
			expected: &proto.DeletePlaylistArtworkResponse{OldArtworkFilename: testImageFilename},
		},
		{
			name: "Error 500. mock.DoesPlaylistExist returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input:       &proto.DeletePlaylistArtworkOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 400. Playlist does not exist",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return false, nil
				},
			},
			input:       &proto.DeletePlaylistArtworkOptions{},
			expectedErr: true,
			err:         status.Error(codes.NotFound, constants.PlaylistNotFoundMessage),
		},
		{
			name: "Error 500. mock.IsOwner returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, errors.New("error")
				},
			},
			input:       &proto.DeletePlaylistArtworkOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 403. User is not playlist owner",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return false, nil
				},
			},
			input:       &proto.DeletePlaylistArtworkOptions{},
			expectedErr: true,
			err:         status.Error(codes.PermissionDenied, constants.NotPlaylistOwnerMessage),
		},
		{
			name: "Error 500. mock.GetOldPlaylistSettings returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return "", errors.New("error")
				},
			},
			input:       &proto.DeletePlaylistArtworkOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
		{
			name: "Error 500. mock.DeletePlaylistArtwork returned error",
			storageMock: &mock.MockStorage{
				DoesPlaylistExistFunc: func(int64) (bool, error) {
					return true, nil
				},
				IsOwnerFunc: func(int64, int64) (bool, error) {
					return true, nil
				},
				GetOldPlaylistSettingsFunc: func(int64) (string, error) {
					return testImageFilename, nil
				},
				DeletePlaylistArtworkFunc: func(int64) error {
					return errors.New("error")
				},
			},
			input:       &proto.DeletePlaylistArtworkOptions{},
			expectedErr: true,
			err:         status.Error(codes.Internal, "error"),
		},
	}

	for _, test := range tests {
		currentTest := test
		t.Run(currentTest.name, func(t *testing.T) {
			storage := NewPlaylistsService(currentTest.storageMock)

			res, err := storage.DeletePlaylistArtwork(context.Background(), currentTest.input)
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

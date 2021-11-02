package usecase

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlaylistUseCase_GetHome(t *testing.T) {
	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}
	tests := []struct {
		name          string
		amount        int
		dbMock        *mock.MockPlayListRepository
		expected      []models.Playlist
		expectedError bool
	}{
		{
			name:   "get 4 playlists",
			amount: 4,
			dbMock: &mock.MockPlayListRepository{
				GetFunc: func(amount int, id int) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < amount; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 4; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 playlists",
			amount: 10,
			dbMock: &mock.MockPlayListRepository{
				GetFunc: func(amount int, id int) ([]models.Playlist, error) {
					var playlists []models.Playlist
					for i := 0; i < amount; i++ {
						playlists = append(playlists, playlist)
					}
					return playlists, nil
				},
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 10; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 error",
			amount: 10,
			dbMock: &mock.MockPlayListRepository{
				GetFunc: func(amount int, id int) ([]models.Playlist, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 10; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewPlaylistUseCase(test.dbMock)
			result, err := useCase.GetHome(test.amount)
			if test.expectedError {
				assert.Error(t, err.OriginalError)
			} else {
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

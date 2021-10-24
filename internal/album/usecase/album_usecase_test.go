package usecase

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAlbumUseCase_GetHome(t *testing.T) {
	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		Artist:         "awa",
		Artwork:        "awa",
		TracksCount:    1,
		TracksDuration: 1,
	}
	tests := []struct {
		name          string
		amount        int
		dbMock        *mock.MockAlbumRepository
		expected      []models.Album
		expectedError bool
	}{
		{
			name:   "get 4 albums",
			amount: 4,
			dbMock: &mock.MockAlbumRepository{
				GetRandomFunc: func(amount int) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < amount; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 4; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 albums",
			amount: 10,
			dbMock: &mock.MockAlbumRepository{
				GetRandomFunc: func(amount int) ([]models.Album, error) {
					var albums []models.Album
					for i := 0; i < amount; i++ {
						albums = append(albums, album)
					}
					return albums, nil
				},
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 10; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: false,
		},
		{
			name:   "get albums error",
			amount: 10,
			dbMock: &mock.MockAlbumRepository{
				GetRandomFunc: func(amount int) ([]models.Album, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 10; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewAlbumUseCase(test.dbMock)
			result, err := useCase.GetHome(test.amount)
			if test.expectedError {
				assert.Error(t, err.OriginalError)
			} else {
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

package usecase

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTrackUseCase_GetHome(t *testing.T) {
	track := models.Track{
		Id:          1,
		Title:       "awa",
		Artist:      "awa",
		Album:       "awa",
		Explicit:    true,
		Genre:       "awa",
		Number:      1,
		File:        "awa",
		ListenCount: 1,
		Duration:    1,
		Lossless:    true,
		Cover:       "awa",
	}

	trackWithoutFile := track
	trackWithoutFile.File = ""

	tests := []struct {
		name          string
		amount        int
		isAuthorized  bool
		dbMock        *mock.MockTrackRepository
		expected      []models.Track
		expectedError bool
	}{
		{
			name:         "get 4 tracks",
			amount:       4,
			isAuthorized: true,
			dbMock: &mock.MockTrackRepository{
				GetRandomFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, track)
					}
					return tracks, nil
				},
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 4; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: false,
		},
		{
			name:         "get 10 tracks",
			amount:       10,
			isAuthorized: true,
			dbMock: &mock.MockTrackRepository{
				GetRandomFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, track)
					}
					return tracks, nil
				},
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 10; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: false,
		},
		{
			name:         "get 10 tracks unauthorized",
			amount:       10,
			isAuthorized: true,
			dbMock: &mock.MockTrackRepository{
				GetRandomFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					var tracks []models.Track
					for i := 0; i < amount; i++ {
						tracks = append(tracks, trackWithoutFile)
					}
					return tracks, nil
				},
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 10; i++ {
					tracks = append(tracks, trackWithoutFile)
				}
				return tracks
			}(),
			expectedError: false,
		},
		{
			name:         "get tracks error",
			amount:       1,
			isAuthorized: true,
			dbMock: &mock.MockTrackRepository{
				GetRandomFunc: func(amount int, isAuthorized bool) ([]models.Track, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 1; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewTrackUseCase(test.dbMock)
			result, err := useCase.GetHome(test.amount, test.isAuthorized)
			if test.expectedError {
				assert.Error(t, err.OriginalError)
			} else {
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

package usecase

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArtistUseCase_GetProfile(t *testing.T) {
	track := models.Track{
		Id:       1,
		Title:    "awa",
		Explicit: true,
		File:     "awa",
		Duration: 1,
		Lossless: true,
		Cover:    "awa",
	}
	tracksUnAuth := track
	tracksUnAuth.File = ""

	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		ArtWork:        "awa",
		TracksDuration: 1,
	}

	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
		Tracks: []models.Track{track},
		Albums: []models.Album{album},
	}
	artistUnAuth := artist
	artistUnAuth.Tracks = []models.Track{tracksUnAuth}

	tests := []struct {
		name          string
		id            int
		isAuthorized  bool
		dbMock        *mock.MockArtistRepository
		expected      *models.Artist
		expectedError bool
	}{
		{
			name: "get profile",
			id:   135,
			dbMock: &mock.MockArtistRepository{
				GetFunc: func(id int) (*models.Artist, error) {
					return &models.Artist{
						Id:     1,
						Name:   "awa",
						Avatar: "awa",
					}, nil
				},
				GetTracksFunc: func(id int, isAuthorized bool, amount int) ([]models.Track, error) {
					return []models.Track{track}, nil
				},
				GetAlbumsFunc: func(id int, amount int) ([]models.Album, error) {
					return []models.Album{album}, nil
				},
			},
			expected:      &artist,
			expectedError: false,
		},
		{
			name: "get profile unauthorized",
			id:   135,
			dbMock: &mock.MockArtistRepository{
				GetFunc: func(id int) (*models.Artist, error) {
					return &models.Artist{
						Id:     1,
						Name:   "awa",
						Avatar: "awa",
					}, nil
				},
				GetTracksFunc: func(id int, isAuthorized bool, amount int) ([]models.Track, error) {
					return []models.Track{tracksUnAuth}, nil
				},
				GetAlbumsFunc: func(id int, amount int) ([]models.Album, error) {
					return []models.Album{album}, nil
				},
			},
			expected:      &artistUnAuth,
			expectedError: false,
		},
		{
			name: "get profile unauthorized",
			id:   135,
			dbMock: &mock.MockArtistRepository{
				GetFunc: func(id int) (*models.Artist, error) {
					return &models.Artist{
						Id:     1,
						Name:   "awa",
						Avatar: "awa",
					}, nil
				},
				GetTracksFunc: func(id int, isAuthorized bool, amount int) ([]models.Track, error) {
					return []models.Track{tracksUnAuth}, nil
				},
				GetAlbumsFunc: func(id int, amount int) ([]models.Album, error) {
					return []models.Album{album}, nil
				},
			},
			expected:      &artistUnAuth,
			expectedError: false,
		},
		{
			name: "Get() error",
			id:   135,
			dbMock: &mock.MockArtistRepository{
				GetFunc: func(id int) (*models.Artist, error) {
					return nil, errors.New("error")
				},
				GetTracksFunc: func(id int, isAuthorized bool, amount int) ([]models.Track, error) {
					return []models.Track{tracksUnAuth}, nil
				},
				GetAlbumsFunc: func(id int, amount int) ([]models.Album, error) {
					return []models.Album{album}, nil
				},
			},
			expected:      &artistUnAuth,
			expectedError: true,
		},
		{
			name: "GetTracks() error",
			id:   135,
			dbMock: &mock.MockArtistRepository{
				GetFunc: func(id int) (*models.Artist, error) {
					return &models.Artist{
						Id:     1,
						Name:   "awa",
						Avatar: "awa",
					}, nil
				},
				GetTracksFunc: func(id int, isAuthorized bool, amount int) ([]models.Track, error) {
					return nil, errors.New("error")
				},
				GetAlbumsFunc: func(id int, amount int) ([]models.Album, error) {
					return []models.Album{album}, nil
				},
			},
			expected:      &artistUnAuth,
			expectedError: true,
		},
		{
			name: "GetAlbums() error",
			id:   135,
			dbMock: &mock.MockArtistRepository{
				GetFunc: func(id int) (*models.Artist, error) {
					return &models.Artist{
						Id:     1,
						Name:   "awa",
						Avatar: "awa",
					}, nil
				},
				GetTracksFunc: func(id int, isAuthorized bool, amount int) ([]models.Track, error) {
					return []models.Track{tracksUnAuth}, nil
				},
				GetAlbumsFunc: func(id int, amount int) ([]models.Album, error) {
					return nil, errors.New("error")
				},
			},
			expected:      &artistUnAuth,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewArtistUseCase(test.dbMock)
			result, err := useCase.GetProfile(test.id, test.isAuthorized)
			if test.expectedError {
				assert.Error(t, err.OriginalError)
			} else {
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

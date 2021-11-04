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
		Id:          1,
		Title:       "awa",
		Explicit:    true,
		Genre:       "awa",
		Number:      1,
		File:        "awa",
		ListenCount: 1,
		Duration:    1,
		Lossless:    true,
		Album: models.Album{
			Id:      1,
			Title:   "awa",
			Artwork: "awa",
		},
	}
	tracksUnAuth := track
	tracksUnAuth.File = ""

	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		Artwork:        "awa",
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

func TestArtistUseCase_GetHome(t *testing.T) {
	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
	}
	tests := []struct {
		name          string
		amount        int
		dbMock        *mock.MockArtistRepository
		expected      []models.Artist
		expectedError bool
	}{
		{
			name:   "get 4 artists",
			amount: 4,
			dbMock: &mock.MockArtistRepository{
				GetRandomFunc: func(amount int) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < amount; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 4; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 artists",
			amount: 10,
			dbMock: &mock.MockArtistRepository{
				GetRandomFunc: func(amount int) ([]models.Artist, error) {
					var artists []models.Artist
					for i := 0; i < amount; i++ {
						artists = append(artists, artist)
					}
					return artists, nil
				},
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 10; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: false,
		},
		{
			name:   "get artists error",
			amount: 10,
			dbMock: &mock.MockArtistRepository{
				GetRandomFunc: func(amount int) ([]models.Artist, error) {
					return nil, errors.New("error")
				},
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 10; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			useCase := NewArtistUseCase(test.dbMock)
			result, err := useCase.GetHome(test.amount)
			if test.expectedError {
				assert.Error(t, err.OriginalError)
			} else {
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
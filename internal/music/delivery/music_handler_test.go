package delivery

import (
	"2021_2_LostPointer/internal/mock"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/music/usecase"
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMusicDelivery_Home(t *testing.T) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	prLogger, _ := config.Build()
	logger := prLogger.Sugar()
	defer prLogger.Sync()

	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}

	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
	}

	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		Artist:         "awa",
		ArtWork:        "awa",
		TracksCount:    1,
		TracksDuration: 1,
	}

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

	collection := models.MusicCollection{
		Tracks: func() []models.Track {
			var tracks []models.Track
			for i := 0; i < usecase.TracksCollectionLimit; i++ {
				tracks = append(tracks, track)
			}
			return tracks
		}(),
		Albums: func() []models.Album {
			var albums []models.Album
			for i := 0; i < usecase.AlbumCollectionLimit; i++ {
				albums = append(albums, album)
			}
			return albums
		}(),
		Artists: func() []models.Artist {
			var artists []models.Artist
			for i := 0; i < usecase.ArtistsCollectionLimit; i++ {
				artists = append(artists, artist)
			}
			return artists
		}(),
		Playlists: func() []models.Playlist {
			var playlists []models.Playlist
			for i := 0; i < usecase.ArtistsCollectionLimit; i++ {
				playlists = append(playlists, playlist)
			}
			return playlists
		}(),
	}

	tests := []struct {
		name          string
		useCaseMock   *mock.MockMusicUseCase
		expected      models.MusicCollection
		expectedError bool
	}{
		{
			name:          "default test",
			useCaseMock:   &mock.MockMusicUseCase{
				GetMusicCollectionFunc: func(isAuthorized bool) (*models.MusicCollection, *models.CustomError) {
					return &collection, nil
				},
			},
			expected:      collection,
			expectedError: false,
		},
		{
			name:          "InternalServeError test",
			useCaseMock:   &mock.MockMusicUseCase{
				GetMusicCollectionFunc: func(isAuthorized bool) (*models.MusicCollection, *models.CustomError	) {
					return nil, &models.CustomError{
						ErrorType: http.StatusInternalServerError,
						OriginalError: errors.New("error"),
						Message: "error",
					}
				},
			},
			expected:      collection,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := echo.New()
			request := httptest.NewRequest(echo.GET, "/", nil)
			recorder := httptest.NewRecorder()
			ctx := server.NewContext(request, recorder)
			ctx.SetPath("api/v1/home")

			ctx.Set("REQUEST_ID", "1")
			ctx.Set("IS_AUTHORIZED", true)

			delivery := NewMusicDelivery(test.useCaseMock, logger)
			_ = delivery.Home(ctx)

			body := recorder.Body
			status := recorder.Result().Status

			var result models.MusicCollection
			_ = json.Unmarshal(body.Bytes(), &result)

			if test.expectedError {
				assert.Equal(t, "500 Internal Server Error", status)
			} else {
				assert.Equal(t, test.expected, result)
				assert.Equal(t, "200 OK", status)
			}
		})
	}
}

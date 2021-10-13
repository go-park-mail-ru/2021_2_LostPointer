package delivery

import (
	"2021_2_LostPointer/pkg/models"
	"2021_2_LostPointer/pkg/music/mock"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestMusicDelivery_Home(t *testing.T) {
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
		ListenCount: 1,
		Duration:    1,
		Cover:       "awa",
	}

	collection := models.MusicCollection{
		Tracks:    []models.Track{track, track, track, track},
		Albums:    []models.Album{album, album, album, album},
		Artists:   []models.Artist{artist, artist, artist, artist},
		Playlists: []models.Playlist{playlist, playlist, playlist, playlist},
	}

	usecaseMock := &mock.MockMusicUseCaseIFace{
		GetMusicCollectionFunc: func() (*models.MusicCollection, error) {
			return &collection, nil
		},
	}

	tests := []struct {
		name          string
		usecaseMock   *mock.MockMusicUseCaseIFace
		expected      models.MusicCollection
		expectedError bool
	}{
		{
			name: "default test",
			usecaseMock: usecaseMock,
			expected: collection,
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := echo.New()
			request := httptest.NewRequest(echo.GET, "/", nil)
			recorder := httptest.NewRecorder()
			ctx := server.NewContext(request, recorder)
			ctx.SetPath("api/v1/home")

			delivery := NewMusicDelivery(test.usecaseMock)
			_ = delivery.Home(ctx)

			body := recorder.Body


			var result models.MusicCollection
			_ = json.Unmarshal(body.Bytes(), &result)


			if test.expectedError {
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

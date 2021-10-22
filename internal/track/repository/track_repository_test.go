package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestTrackRepository_GetRandom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewTrackRepository(db)

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
		mock          func()
		expected      []models.Track
		expectedError bool
	}{
		{
			name:         "get 4 random tracks",
			amount:       4,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "cover"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Cover)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, art.name, alb.title, explicit, " +
					"g.name, number, file, listen_count, duration, lossless, alb.artwork as cover FROM tracks " +
					"LEFT JOIN genres g ON tracks.genre = g.id " +
					"LEFT JOIN albums alb ON tracks.album = alb.id " +
					"LEFT JOIN artists art ON tracks.artist = art.id ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(4)).WillReturnRows(rows)
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
			name:         "get 10 random tracks",
			amount:       10,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "cover"})
				for i := 0; i < 10; i++ {
					rows.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Cover)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, art.name, alb.title, explicit, " +
					"g.name, number, file, listen_count, duration, lossless, alb.artwork as cover FROM tracks " +
					"LEFT JOIN genres g ON tracks.genre = g.id " +
					"LEFT JOIN albums alb ON tracks.album = alb.id " +
					"LEFT JOIN artists art ON tracks.artist = art.id ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(10)).WillReturnRows(rows)
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
			name:         "get 100 random tracks",
			amount:       100,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "cover"})
				for i := 0; i < 100; i++ {
					rows.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Cover)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, art.name, alb.title, explicit, " +
					"g.name, number, file, listen_count, duration, lossless, alb.artwork as cover FROM tracks " +
					"LEFT JOIN genres g ON tracks.genre = g.id " +
					"LEFT JOIN albums alb ON tracks.album = alb.id " +
					"LEFT JOIN artists art ON tracks.artist = art.id ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() []models.Track {
				var tracks []models.Track
				for i := 0; i < 100; i++ {
					tracks = append(tracks, track)
				}
				return tracks
			}(),
			expectedError: false,
		},
		{
			name:         "query error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "cover"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Cover)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, art.name, alb.title, explicit, " +
					"g.name, number, file, listen_count, duration, lossless, alb.artwork as cover FROM tracks " +
					"LEFT JOIN genres g ON tracks.genre = g.id " +
					"LEFT JOIN albums alb ON tracks.album = alb.id " +
					"LEFT JOIN artists art ON tracks.artist = art.id ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(1)).WillReturnError(errors.New("error"))
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
		{
			name:         "scan error",
			amount:       1,
			isAuthorized: true,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "cover", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Cover, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, art.name, alb.title, explicit, " +
					"g.name, number, file, listen_count, duration, lossless, alb.artwork as cover FROM tracks " +
					"LEFT JOIN genres g ON tracks.genre = g.id " +
					"LEFT JOIN albums alb ON tracks.album = alb.id " +
					"LEFT JOIN artists art ON tracks.artist = art.id ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(1)).WillReturnRows(rows)
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
		{
			name:         "get 4 random tracks unauthorized",
			amount:       4,
			isAuthorized: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "cover"})
				for i := 0; i < 4; i++ {
					rows.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, "", track.ListenCount, track.Duration, track.Lossless, track.Cover)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, art.name, alb.title, explicit, " +
					"g.name, number, file, listen_count, duration, lossless, alb.artwork as cover FROM tracks " +
					"LEFT JOIN genres g ON tracks.genre = g.id " +
					"LEFT JOIN albums alb ON tracks.album = alb.id " +
					"LEFT JOIN artists art ON tracks.artist = art.id ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() []models.Track {
				var tracks []models.Track

				for i := 0; i < 4; i++ {
					tracks = append(tracks, trackWithoutFile)
				}
				return tracks
			}(),
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetRandom(test.amount, test.isAuthorized)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

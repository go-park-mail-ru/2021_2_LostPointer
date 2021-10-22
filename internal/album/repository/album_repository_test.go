package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestAlbumRepository_GetRandom(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewAlbumRepository(db)

	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		Artist:         "awa",
		ArtWork:        "awa",
		TracksCount:    1,
		TracksDuration: 1,
	}

	tests := []struct {
		name          string
		amount        int
		mock          func()
		expected      []models.Album
		expectedError bool
	}{
		{
			name:   "get 4 random albums",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "art.name", "a.artwork", "a.track_count", "tracksDuration"})
				for i := 0; i < 4; i++ {
					rows.AddRow(album.Id, album.Title, album.Year, album.Artist, album.ArtWork, album.TracksCount, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.year, art.name, " +
					"a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a " +
					"LEFT JOIN artists art ON art.id = a.artist " +
					"JOIN tracks t on t.album = a.id " +
					"WHERE art.name NOT LIKE '%Frank Sinatra%' " +
					"GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count " +
					"ORDER BY RANDOM() " +
					"LIMIT $1")).WithArgs(driver.Value(4)).WillReturnRows(rows)
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
			name:   "get 10 random albums",
			amount: 10,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "art.name", "a.artwork", "a.track_count", "tracksDuration"})
				for i := 0; i < 10; i++ {
					rows.AddRow(album.Id, album.Title, album.Year, album.Artist, album.ArtWork, album.TracksCount, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.year, art.name, " +
					"a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a " +
					"LEFT JOIN artists art ON art.id = a.artist " +
					"JOIN tracks t on t.album = a.id " +
					"WHERE art.name NOT LIKE '%Frank Sinatra%' " +
					"GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count " +
					"ORDER BY RANDOM() " +
					"LIMIT $1")).WithArgs(driver.Value(10)).WillReturnRows(rows)
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
			name:   "get 100 random albums",
			amount: 100,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "art.name", "a.artwork", "a.track_count", "tracksDuration"})
				for i := 0; i < 100; i++ {
					rows.AddRow(album.Id, album.Title, album.Year, album.Artist, album.ArtWork, album.TracksCount, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.year, art.name, " +
					"a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a " +
					"LEFT JOIN artists art ON art.id = a.artist " +
					"JOIN tracks t on t.album = a.id " +
					"WHERE art.name NOT LIKE '%Frank Sinatra%' " +
					"GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count " +
					"ORDER BY RANDOM() " +
					"LIMIT $1")).WithArgs(driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 100; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: false,
		},
		{
			name:   "query error",
			amount: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "art.name", "a.artwork", "a.track_count", "tracksDuration"})
				for i := 0; i < 1; i++ {
					rows.AddRow(album.Id, album.Title, album.Year, album.Artist, album.ArtWork, album.TracksCount, album.TracksDuration)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.year, art.name, " +
					"a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a " +
					"LEFT JOIN artists art ON art.id = a.artist " +
					"JOIN tracks t on t.album = a.id " +
					"WHERE art.name NOT LIKE '%Frank Sinatra%' " +
					"GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count " +
					"ORDER BY RANDOM() " +
					"LIMIT $1")).WithArgs(driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 1; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
		{
			name:   "scan error",
			amount: 1,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "art.name", "a.artwork", "a.track_count", "tracksDuration", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(album.Id, album.Title, album.Year, album.Artist, album.ArtWork, album.TracksCount, album.TracksDuration, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.year, art.name, " +
					"a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a " +
					"LEFT JOIN artists art ON art.id = a.artist " +
					"JOIN tracks t on t.album = a.id " +
					"WHERE art.name NOT LIKE '%Frank Sinatra%' " +
					"GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count " +
					"ORDER BY RANDOM() " +
					"LIMIT $1")).WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() []models.Album {
				var albums []models.Album
				for i := 0; i < 1; i++ {
					albums = append(albums, album)
				}
				return albums
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetRandom(test.amount)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

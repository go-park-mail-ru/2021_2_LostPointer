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

func TestArtistRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewArtistRepository(db)

	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
	}

	tests := []struct {
		name          string
		id            int
		mock          func()
		expected      models.Artist
		expectedError bool
	}{
		{
			name: "get artist id 135",
			id:   135,
			mock: func() {
				rows := sqlmock.NewRows([]string{"art.id", "art.name", "art.avatar"})
				rows.AddRow(artist.Id, artist.Name, artist.Avatar)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT art.id, art.name, art.avatar FROM artists art " +
					"WHERE art.id = $1 " +
					"GROUP BY art.id")).WithArgs(driver.Value(135)).WillReturnRows(rows)
			},
			expected:      artist,
			expectedError: false,
		},
		{
			name: "database error",
			id:   135,
			mock: func() {
				rows := sqlmock.NewRows([]string{"art.id", "art.name", "art.avatar"})
				rows.AddRow(artist.Id, artist.Name, artist.Avatar)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT art.id, art.name, art.avatar FROM artists art " +
					"WHERE art.id = $1 " +
					"GROUP BY art.id")).WithArgs(driver.Value(135)).WillReturnError(errors.New("error"))
			},
			expected:      artist,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.Get(test.id)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, *result)
			}
		})
	}
}

func TestArtistRepository_GetTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewArtistRepository(db)

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

	tests := []struct {
		name          string
		id            int
		isAuthorized  bool
		mock          func()
		expected      []models.Track
		expectedError bool
	}{
		{
			name:         "get tracks with artist id 135",
			id:           135,
			isAuthorized: true,
			mock: func() {
				rows := sqlmock.NewRows([]string{"t.id", "t.title", "explicit", "file", "duration", "lossless", "alb.artwork"})
				rows.AddRow(track.Id, track.Title, track.Explicit, track.File, track.Duration, track.Lossless, track.Cover)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT t.id, t.title, explicit, file, duration, lossless, " +
					"alb.artwork FROM tracks t " +
					"LEFT JOIN albums alb ON alb.id = t.album " +
					"WHERE t.artist = $1 " +
					"ORDER BY t.listen_count LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnRows(rows)
			},
			expected:      []models.Track{track},
			expectedError: false,
		},
		{
			name:         "get tracks with artist id  unauthorized",
			id:           135,
			isAuthorized: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"t.id", "t.title", "explicit", "file", "duration", "lossless", "alb.artwork"})
				rows.AddRow(track.Id, track.Title, track.Explicit, track.File, track.Duration, track.Lossless, track.Cover)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT t.id, t.title, explicit, file, duration, lossless, " +
					"alb.artwork FROM tracks t " +
					"LEFT JOIN albums alb ON alb.id = t.album " +
					"WHERE t.artist = $1 " +
					"ORDER BY t.listen_count LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnRows(rows)
			},
			expected:      []models.Track{tracksUnAuth},
			expectedError: false,
		},
		{
			name:         "query error",
			id:           135,
			isAuthorized: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"t.id", "t.title", "explicit", "file", "duration", "lossless", "alb.artwork"})
				rows.AddRow(track.Id, track.Title, track.Explicit, track.File, track.Duration, track.Lossless, track.Cover)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT t.id, t.title, explicit, file, duration, lossless, " +
					"alb.artwork FROM tracks t " +
					"LEFT JOIN albums alb ON alb.id = t.album " +
					"WHERE t.artist = $1 " +
					"ORDER BY t.listen_count LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected:      []models.Track{tracksUnAuth},
			expectedError: true,
		},
		{
			name:         "scan error",
			id:           135,
			isAuthorized: false,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"t.id", "t.title", "explicit", "file", "duration", "lossless", "alb.artwork", "newArg"})
				rows.AddRow(track.Id, track.Title, track.Explicit, track.File, track.Duration, track.Lossless, track.Cover, newArg)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT t.id, t.title, explicit, file, duration, lossless, " +
					"alb.artwork FROM tracks t " +
					"LEFT JOIN albums alb ON alb.id = t.album " +
					"WHERE t.artist = $1 " +
					"ORDER BY t.listen_count LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected:      []models.Track{tracksUnAuth},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetTracks(test.id, test.isAuthorized, 1)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestArtistRepository_GetAlbums(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewArtistRepository(db)

	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           1,
		ArtWork:        "awa",
		TracksDuration: 1,
	}

	tests := []struct {
		name          string
		id            int
		mock          func()
		expected      []models.Album
		expectedError bool
	}{
		{
			name: "get albums with artist id 135",
			id:   135,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.artwork", "a.year", "tracksDuration"})
				rows.AddRow(album.Id, album.Title, album.ArtWork, album.Year, album.TracksDuration)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.artwork, a.year, SUM(t.duration) AS " +
					"tracksDuration FROM albums a " +
					"JOIN tracks t on t.album = a.id " +
					"WHERE a.artist = $1 " +
					"GROUP BY a.id, a.title, a.artwork, a.year " +
					"ORDER BY a.year DESC")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnRows(rows)
			},
			expected:      []models.Album{album},
			expectedError: false,
		},
		{
			name: "query error",
			id:   135,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.artwork", "a.year", "tracksDuration"})
				rows.AddRow(album.Id, album.Title, album.ArtWork, album.Year, album.TracksDuration)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.artwork, a.year, SUM(t.duration) AS " +
					"tracksDuration FROM albums a " +
					"JOIN tracks t on t.album = a.id " +
					"WHERE a.artist = $1 " +
					"GROUP BY a.id, a.title, a.artwork, a.year " +
					"ORDER BY a.year DESC")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected:      []models.Album{album},
			expectedError: true,
		},
		{
			name: "scan error",
			id:   135,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.artwork", "a.year", "tracksDuration", "newArg"})
				rows.AddRow(album.Id, album.Title, album.ArtWork, album.Year, album.TracksDuration, newArg)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.artwork, a.year, SUM(t.duration) AS " +
					"tracksDuration FROM albums a " +
					"JOIN tracks t on t.album = a.id " +
					"WHERE a.artist = $1 " +
					"GROUP BY a.id, a.title, a.artwork, a.year " +
					"ORDER BY a.year DESC")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected:      []models.Album{album},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetAlbums(test.id, 1)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

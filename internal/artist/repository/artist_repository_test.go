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
		Avatar: "avatar.webp",
		Video:  "video.mp4",
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
				rows := sqlmock.NewRows([]string{"art.id", "art.name", "art.avatar", "art.video"})
				rows.AddRow(artist.Id, artist.Name, "avatar", "video")
				mock.ExpectQuery(regexp.QuoteMeta("SELECT art.id, art.name, art.avatar, art.video FROM artists art " +
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
				rows := sqlmock.NewRows([]string{"art.id", "art.name", "art.avatar", "art.video"})
				rows.AddRow(artist.Id, artist.Name, "avatar", "video")
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
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork"})
				rows.AddRow(track.Id, track.Title, track.Explicit, track.Genre, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.Id, track.Album.Title, track.Album.Artwork)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, explicit, "+
					"g.name, number, file, listen_count, duration, lossless, alb.id, alb.title, alb.artwork as cover FROM tracks "+
					"LEFT JOIN genres g ON tracks.genre = g.id "+
					"LEFT JOIN albums alb ON tracks.album = alb.id "+
					"LEFT JOIN artists art ON tracks.artist = art.id "+
					"WHERE tracks.artist = $1 "+
					"ORDER BY tracks.listen_count DESC LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnRows(rows)
			},
			expected:      []models.Track{track},
			expectedError: false,
		},
		{
			name:         "get tracks with artist id 135 unauthorized",
			id:           135,
			isAuthorized: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork"})
				rows.AddRow(track.Id, track.Title, track.Explicit, track.Genre, track.Number, "", track.ListenCount, track.Duration, track.Lossless, track.Album.Id, track.Album.Title, track.Album.Artwork)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, explicit, "+
					"g.name, number, file, listen_count, duration, lossless, alb.id, alb.title, alb.artwork as cover FROM tracks "+
					"LEFT JOIN genres g ON tracks.genre = g.id "+
					"LEFT JOIN albums alb ON tracks.album = alb.id "+
					"LEFT JOIN artists art ON tracks.artist = art.id "+
					"WHERE tracks.artist = $1 "+
					"ORDER BY tracks.listen_count DESC LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnRows(rows)
			},
			expected:      []models.Track{tracksUnAuth},
			expectedError: false,
		},
		{
			name:         "query error",
			id:           135,
			isAuthorized: false,
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork"})
				rows.AddRow(track.Id, track.Title, track.Explicit, track.Genre, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.Id, track.Album.Title, track.Album.Artwork)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, explicit, "+
					"g.name, number, file, listen_count, duration, lossless, alb.id, alb.title, alb.artwork as cover FROM tracks "+
					"LEFT JOIN genres g ON tracks.genre = g.id "+
					"LEFT JOIN albums alb ON tracks.album = alb.id "+
					"LEFT JOIN artists art ON tracks.artist = art.id "+
					"WHERE tracks.artist = $1 "+
					"ORDER BY tracks.listen_count DESC LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnError(errors.New("error"))
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
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "explicit", "g.name", "number", "file", "listen_count", "duration", "lossless", "alb.id", "alb.title", "alb.artwork", "newArg"})
				rows.AddRow(track.Id, track.Title, track.Explicit, track.Genre, track.Number, track.File, track.ListenCount, track.Duration, track.Lossless, track.Album.Id, track.Album.Title, track.Album.Artwork, newArg)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT tracks.id, tracks.title, explicit, "+
					"g.name, number, file, listen_count, duration, lossless, alb.id, alb.title, alb.artwork as cover FROM tracks "+
					"LEFT JOIN genres g ON tracks.genre = g.id "+
					"LEFT JOIN albums alb ON tracks.album = alb.id "+
					"LEFT JOIN artists art ON tracks.artist = art.id "+
					"WHERE tracks.artist = $1 "+
					"ORDER BY tracks.listen_count DESC LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnRows(rows)
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
		Artwork:        "awa",
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
				rows.AddRow(album.Id, album.Title, album.Artwork, album.Year, album.TracksDuration)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.artwork, a.year, SUM(t.duration) AS "+
					"tracksDuration FROM albums a "+
					"JOIN tracks t on t.album = a.id "+
					"WHERE a.artist = $1 "+
					"GROUP BY a.id, a.title, a.artwork, a.year "+
					"ORDER BY a.year DESC LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnRows(rows)
			},
			expected:      []models.Album{album},
			expectedError: false,
		},
		{
			name: "query error",
			id:   135,
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.artwork", "a.year", "tracksDuration"})
				rows.AddRow(album.Id, album.Title, album.Artwork, album.Year, album.TracksDuration)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.artwork, a.year, SUM(t.duration) AS "+
					"tracksDuration FROM albums a "+
					"JOIN tracks t on t.album = a.id "+
					"WHERE a.artist = $1 "+
					"GROUP BY a.id, a.title, a.artwork, a.year "+
					"ORDER BY a.year DESC LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnError(errors.New("error"))
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
				rows.AddRow(album.Id, album.Title, album.Artwork, album.Year, album.TracksDuration, newArg)
				mock.ExpectQuery(regexp.QuoteMeta("SELECT a.id, a.title, a.artwork, a.year, SUM(t.duration) AS "+
					"tracksDuration FROM albums a "+
					"JOIN tracks t on t.album = a.id "+
					"WHERE a.artist = $1 "+
					"GROUP BY a.id, a.title, a.artwork, a.year "+
					"ORDER BY a.year DESC LIMIT $2")).WithArgs(driver.Value(135), driver.Value(1)).WillReturnRows(rows)
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

func TestArtistRepository_GetRandom(t *testing.T) {
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
		amount        int
		mock          func()
		expected      []models.Artist
		expectedError bool
	}{
		{
			name:   "get 4 random artists",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < 4; i++ {
					rows.AddRow(artist.Id, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT artists.id, artists.name, artists.avatar FROM artists " +
					"WHERE artists.name NOT LIKE '%Frank Sinatra%' ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(4)).WillReturnRows(rows)
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
			name:   "get 10 random artists",
			amount: 10,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < 10; i++ {
					rows.AddRow(artist.Id, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT artists.id, artists.name, artists.avatar FROM artists " +
					"WHERE artists.name NOT LIKE '%Frank Sinatra%' ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(10)).WillReturnRows(rows)
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
			name:   "get 100 random artists",
			amount: 100,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < 100; i++ {
					rows.AddRow(artist.Id, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT artists.id, artists.name, artists.avatar FROM artists " +
					"WHERE artists.name NOT LIKE '%Frank Sinatra%' ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 100; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: false,
		},
		{
			name:   "query error",
			amount: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				for i := 0; i < 1; i++ {
					rows.AddRow(artist.Id, artist.Name, artist.Avatar)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT artists.id, artists.name, artists.avatar FROM artists " +
					"WHERE artists.name NOT LIKE '%Frank Sinatra%' ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 1; i++ {
					artists = append(artists, artist)
				}
				return artists
			}(),
			expectedError: true,
		},
		{
			name:   "scan error",
			amount: 1,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(artist.Id, artist.Name, artist.Avatar, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT artists.id, artists.name, artists.avatar FROM artists " +
					"WHERE artists.name NOT LIKE '%Frank Sinatra%' ORDER BY RANDOM() LIMIT $1")).WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() []models.Artist {
				var artists []models.Artist
				for i := 0; i < 1; i++ {
					artists = append(artists, artist)
				}
				return artists
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

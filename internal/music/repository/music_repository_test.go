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

func TestMusicRepository_GetRandomTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewMusicRepository(db)

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
			result, err := repository.GetRandomTracks(test.amount, test.isAuthorized)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicRepository_GetRandomAlbums(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewMusicRepository(db)

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
			result, err := repository.GetRandomAlbums(test.amount)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}


func TestMusicRepository_GetRandomPlaylists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewMusicRepository(db)

	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}

	tests := []struct {
		name          string
		amount        int
		mock          func()
		expected      []models.Playlist
		expectedError bool
	}{
		{
			name:   "get 4 random playlists",
			amount: 4,
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				for i := 0; i < 4; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user " +
					"FROM playlists LIMIT $1")).WithArgs(driver.Value(4)).WillReturnRows(rows)
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 4; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: false,
		},
		{
			name:   "get 10 random playlists",
			amount: 10,
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				for i := 0; i < 10; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user " +
					"FROM playlists LIMIT $1")).WithArgs(driver.Value(10)).WillReturnRows(rows)
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 10; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: false,
		},
		{
			name:   "get 100 random playlists",
			amount: 100,
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				for i := 0; i < 100; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user " +
					"FROM playlists LIMIT $1")).WithArgs(driver.Value(100)).WillReturnRows(rows)
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 100; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: false,
		},
		{
			name:   "query error",
			amount: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				for i := 0; i < 1; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user " +
					"FROM playlists LIMIT $1")).WithArgs(driver.Value(1)).WillReturnError(errors.New("error"))
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 1; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: true,
		},
		{
			name:   "scan error",
			amount: 1,
			mock: func() {
				var newArg = 1
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user", "newArg"})
				for i := 0; i < 1; i++ {
					rows.AddRow(playlist.Id, playlist.Name, playlist.User, newArg)
				}
				mock.ExpectQuery(regexp.QuoteMeta("SELECT playlists.id, playlists.title, playlists.user " +
					"FROM playlists LIMIT $1")).WithArgs(driver.Value(1)).WillReturnRows(rows)
			},
			expected: func() []models.Playlist {
				var playlists []models.Playlist
				for i := 0; i < 1; i++ {
					playlists = append(playlists, playlist)
				}
				return playlists
			}(),
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetRandomPlaylists(test.amount)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

package repository

import (
	"2021_2_LostPointer/pkg/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestMusicRepository_CreateTracksRequestWithParameters(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repository := NewMusicRepository(db)

	tests := []struct {
		name          string
		gettingWith   uint8
		parameters    []string
		distinctOn    uint8
		expected      string
		expectedError bool
	}{
		{
			name:        "GettingWithID; 1, 2, 3, 4, 5; DistinctOnArtists",
			gettingWith: GettingWithID,
			parameters:  []string{"1", "2", "3", "4", "5"},
			distinctOn:  DistinctOnArtists,
			expected: `SELECT DISTINCT ON(art.name) tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE tracks.id IN (1, 2, 3, 4, 5)`,
			expectedError: false,
		},
		{
			name:        "GettingWithID; 1, 2, 3, 4, 5, 6; DistinctOnNone",
			gettingWith: GettingWithID,
			parameters:  []string{"1", "2", "3", "4", "5", "6"},
			distinctOn:  DistinctOnNone,
			expected: `SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE tracks.id IN (1, 2, 3, 4, 5, 6)`,
			expectedError: false,
		},
		{
			name:        "GettingWithID; 1, 2, 3, 4, 5; DistinctOnAlbums",
			gettingWith: GettingWithID,
			parameters:  []string{"1", "2", "3", "4", "5", "6", "7"},
			distinctOn:  DistinctOnAlbums,
			expected: `SELECT DISTINCT ON(alb.title) tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE tracks.id IN (1, 2, 3, 4, 5, 6, 7)`,
			expectedError: false,
		},
		{
			name:        "GettingWithGenres; 'Pop', 'Jazz', 'Easy Listening'; DistinctOnAlbums",
			gettingWith: GettingWithGenres,
			parameters:  []string{"Pop", "Jazz", "Easy Listening"},
			distinctOn:  DistinctOnAlbums,
			expected: `SELECT DISTINCT ON(alb.title) tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop', 'Jazz', 'Easy Listening')`,
			expectedError: false,
		},
		{
			name:        "GettingWithGenres; 'Pop', 'Jazz', 'Easy Listening'; DistinctOnAlbums",
			gettingWith: GettingWithGenres,
			parameters:  []string{"Pop", "Jazz", "Easy Listening"},
			distinctOn:  DistinctOnAlbums,
			expected: `SELECT DISTINCT ON(alb.title) tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop', 'Jazz', 'Easy Listening')`,
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := repository.CreateTracksRequestWithParameters(test.gettingWith, test.parameters, test.distinctOn)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicRepository_CreateAlbumsDefaultRequest(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewMusicRepository(db)

	expected := `SELECT a.id, a.title, a.year, art.name,
									a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a
									LEFT JOIN artists art ON art.id = a.artist
									JOIN tracks t on t.album = a.id
									WHERE art.name = 'Земфира'
									GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
									LIMIT 1`
	result := repository.CreateAlbumsDefaultRequest(1)

	assert.Equal(t, expected, result)
}

func TestMusicRepository_CreateArtistsDefaultRequest(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewMusicRepository(db)

	expected := `SELECT artists.id, artists.name, artists.avatar FROM artists
									WHERE artists.avatar != 'frank_sinatra.jpg' LIMIT 1`
	result := repository.CreateArtistsDefaultRequest(1)

	assert.Equal(t, expected, result)
}

func TestMusicRepository_CreatePlaylistsDefaultRequest(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewMusicRepository(db)

	expected := `SELECT playlists.id, playlists.title, playlists.user FROM playlists LIMIT 1`
	result := repository.CreatePlaylistsDefaultRequest(1)

	assert.Equal(t, expected, result)
}

func TestMusicRepository_GetTracks(t *testing.T) {
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
		ListenCount: 1,
		Duration:    1,
		Cover:       "awa",
	}

	tests := []struct {
		name          string
		mock          func()
		request       func() string
		expected      []models.Track
		expectedError bool
	}{
		{
			name: "GettingWithID; 1; DistinctOnNone",
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration"})
				rows.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, track.Cover, track.ListenCount, track.Duration)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE tracks.id IN (1)`)).WillReturnRows(rows)
			},
			request: func() string {
				request, _ := repository.CreateTracksRequestWithParameters(GettingWithID, []string{"1"}, DistinctOnNone)
				return request
			},
			expected:      []models.Track{track},
			expectedError: false,
		},
		{
			name: "GettingWithGenres; 'Pop'; DistinctOnNone",
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.id", "tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration"})
				rows.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, track.Cover, track.ListenCount, track.Duration)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT tracks.id, tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop')`)).WillReturnRows(rows)
			},
			request: func() string {
				request, _ := repository.CreateTracksRequestWithParameters(GettingWithGenres, []string{"Pop"}, DistinctOnNone)
				return request
			},
			expected:      []models.Track{track},
			expectedError: false,
		},
		{
			name: "Wrong amount of select attributes",
			mock: func() {
				rows := sqlmock.NewRows([]string{"tracks.title", "art.name", "alb.title", "explicit", "g.name", "number", "file", "listen_count", "duration"})
				rows.AddRow(track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number, track.Cover, track.ListenCount, track.Duration)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop')`)).WillReturnRows(rows)
			},
			request: func() string {
				request := `SELECT tracks.title, art.name, alb.title, explicit, g.name, number, file, listen_count, duration FROM tracks
					LEFT JOIN genres g ON tracks.genre = g.id
					LEFT JOIN albums alb ON tracks.album = alb.id
					LEFT JOIN artists art ON tracks.artist = art.id WHERE g.name IN ('Pop')`
				return request
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetTracks(test.request())
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicRepository_GetAlbums(t *testing.T) {
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
		mock          func()
		request       func() string
		expected      []models.Album
		expectedError bool
	}{
		{
			name: "Default request",
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.id", "a.title", "a.year", "art.name", "a.artwork", "a.track_count", "tracksDuration"})
				rows.AddRow(album.Id, album.Title, album.Year, album.Artist, album.ArtWork, album.TracksCount, album.TracksDuration)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.id, a.title, a.year, art.name,
									a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a
									LEFT JOIN artists art ON art.id = a.artist
									JOIN tracks t on t.album = a.id
									WHERE art.name = 'Земфира'
									GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
									LIMIT 1`)).WillReturnRows(rows)
			},
			request: func() string {
				return repository.CreateAlbumsDefaultRequest(1)
			},
			expected:      []models.Album{album},
			expectedError: false,
		},
		{
			name: "Wrong amount of select attributes",
			mock: func() {
				rows := sqlmock.NewRows([]string{"a.title", "a.year", "art.name", "a.artwork", "a.track_count", "tracksDuration"})
				rows.AddRow(album.Title, album.Year, album.Artist, album.ArtWork, album.TracksCount, album.TracksDuration)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT a.title, a.year, art.name,
									a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a
									LEFT JOIN artists art ON art.id = a.artist
									JOIN tracks t on t.album = a.id
									WHERE art.name = 'Земфира'
									GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
									LIMIT 1`)).WillReturnRows(rows)
			},
			request: func() string {
				request := `SELECT a.title, a.year, art.name,
									a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a
									LEFT JOIN artists art ON art.id = a.artist
									JOIN tracks t on t.album = a.id
									WHERE art.name = 'Земфира'
									GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
									LIMIT 1`
				return request
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetAlbums(test.request())
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicRepository_GetArtists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewMusicRepository(db)

	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
	}

	tests := []struct {
		name          string
		mock          func()
		request       func() string
		expected      []models.Artist
		expectedError bool
	}{
		{
			name: "Default request",
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.id", "artists.name", "artists.avatar"})
				rows.AddRow(artist.Id, artist.Name, artist.Avatar)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT artists.id, artists.name, artists.avatar FROM artists
									WHERE artists.avatar != 'frank_sinatra.jpg' LIMIT 1`)).WillReturnRows(rows)
			},
			request: func() string {
				return repository.CreateArtistsDefaultRequest(1)
			},
			expected:      []models.Artist{artist},
			expectedError: false,
		},
		{
			name: "Wrong amount of select attributes",
			mock: func() {
				rows := sqlmock.NewRows([]string{"artists.name", "artists.avatar"})
				rows.AddRow(artist.Name, artist.Avatar)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT artists.name, artists.avatar FROM artists
									WHERE artists.avatar != 'frank_sinatra.jpg' LIMIT 1`)).WillReturnRows(rows)
			},
			request: func() string {
				request := `SELECT artists.name, artists.avatar FROM artists
									WHERE artists.avatar != 'frank_sinatra.jpg' LIMIT 1`
				return request
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetArtists(test.request())
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicRepository_GetPlaylists(t *testing.T) {
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
		mock          func()
		request       func() string
		expected      []models.Playlist
		expectedError bool
	}{
		{
			name: "Default request",
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.id", "playlists.title", "playlists.user"})
				rows.AddRow(playlist.Id, playlist.Name, playlist.User)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT playlists.id, playlists.title, playlists.user FROM playlists LIMIT 1`)).WillReturnRows(rows)
			},
			request: func() string {
				return repository.CreatePlaylistsDefaultRequest(1)
			},
			expected:      []models.Playlist{playlist},
			expectedError: false,
		},
		{
			name: "Wrong amount of select attributes",
			mock: func() {
				rows := sqlmock.NewRows([]string{"playlists.title", "playlists.user"})
				rows.AddRow(playlist.Name, playlist.User)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT playlists.title, playlists.user FROM playlists LIMIT 1`)).WillReturnRows(rows)
			},
			request: func() string {
				request := `SELECT playlists.title, playlists.user FROM playlists LIMIT 1`
				return request
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			result, err := repository.GetPlaylists(test.request())
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

func TestMusicRepository_IsGenreExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repository := NewMusicRepository(db)

	genres := []string{"Pop", "Jazz", "Rap", "Hip-Hop", "Rock"}
	mockFunc := func() {
		rows := sqlmock.NewRows([]string{"name"})
		for _, genre := range genres {
			rows.AddRow(genre)
		}
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT name FROM genres`)).WillReturnRows(rows)
	}

	tests := []struct {
		name          string
		genres        []string
		expected      bool
		expectedError bool
	}{
		{
			name:          "Genre exist",
			genres:        []string{"Pop", "Rap"},
			expected:      true,
			expectedError: false,
		},
		{
			name:          "Genre does not exist",
			genres:        []string{"Pop", "Classic"},
			expected:      false,
			expectedError: false,
		},
		{
			name:          "Empty array",
			genres:        []string{},
			expected:      false,
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockFunc()
			result, err := repository.IsGenreExist(test.genres)
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}

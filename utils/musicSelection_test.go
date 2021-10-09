package utils

import (
	"2021_2_LostPointer/models"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestCreateTracksRequestWithGenres(t *testing.T) {
	result := createTracksRequestWithGenres([]string{"Pop", "Jazz"})
	expected := `SELECT DISTINCT ON(tracks.album) tracks.id FROM tracks
LEFT JOIN genres g on tracks.genre = g.id
WHERE g.name='Pop' OR g.name='Jazz'`

	assert.Equal(t, expected, result)
}

func TestCreateTracksRequestWithIds(t *testing.T) {
	result := createTracksRequestWithIds([]int{1799, 1800})
	expected := `SELECT t.id,
							t.title,
							art.name,
							alb.title,
							t.explicit,
							g.name,
							t.number,
							t.file,
							t.listen_count,
							t.duration,
							alb.artwork
							FROM tracks t
							LEFT JOIN albums alb ON alb.id = t.album
							LEFT JOIN artists art ON art.id = t.artist
							LEFT JOIN genres g ON g.id = t.genre
							WHERE t.id IN (1799, 1800)`

	assert.Equal(t, expected, result)
}

func TestGetTracksIdGenre(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
		SELECT name FROM genres
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"name"})
		row.AddRow("Pop")
		return row
	}())

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
		SELECT DISTINCT ON(tracks.album) tracks.id FROM tracks
LEFT JOIN genres g on tracks.genre = g.id
WHERE g.name='Pop'
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"id"})
		row.AddRow(1799)
		return row
	}())

	result, err := getTracksIdByGenre(db, []string{"Pop"})
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %d\n", err, 1799)
	}

	assert.Equal(t, []int{1799}, result)
}

func TestIsGenreExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
		SELECT name FROM genres
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"name"})
		row.AddRow("Pop")
		return row
	}())

	result, err := isGenresExist([]string{"Pop"}, db)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %d\n", err, 1799)
	}

	assert.Equal(t, true, result)
}

func TestGetTracks(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	track := models.Track{
		Id:          1,
		Title:       "awa",
		Artist:      "awa",
		Album:       "awa",
		Explicit:    false,
		Genre:       "awa",
		Number:      1,
		File:        "awa",
		ListenCount: 1,
		Duration:    1,
		Cover:       "awa",
	}
	var tracks []models.Track
	tracks = append(tracks, track)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
<<<<<<< HEAD
		SELECT name FROM genres
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"name"})
		row.AddRow("Pop")
		return row
	}())

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
		SELECT DISTINCT ON(tracks.album) tracks.id FROM tracks
LEFT JOIN genres g on tracks.genre = g.id
WHERE g.name='Pop'
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"id"})
		row.AddRow(1799)
		row.AddRow(1800)
		return row
	}())

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
		SELECT t.id,
							t.title,
							art.name,
							alb.title,
							t.explicit,
							g.name,
							t.number,
							t.file,
							t.listen_count,
							t.duration,
							alb.artwork
							FROM tracks t
							LEFT JOIN albums alb ON alb.id = t.album
							LEFT JOIN artists art ON art.id = t.artist
							LEFT JOIN genres g ON g.id = t.genre
							WHERE t.id IN (1799)
=======
		SELECT DISTINCT ON (alb.title) t.id,
												                              t.title,
												                              art.name,
												                              alb.title,
												                              t.explicit,
												                              g.name,
												                              t.number,
												                              t.file,
												                              t.listen_count,
												                              t.duration,
												                              alb.artwork
												FROM tracks t
												         LEFT JOIN albums alb ON alb.id = t.album
												         LEFT JOIN artists art ON art.id = t.artist
												         LEFT JOIN genres g ON g.id = t.genre
												WHERE RANDOM() < 0.5
												LIMIT 1
>>>>>>> for_potikat
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"id", "title", "artist", "album", "explicit", "genre", "number",
			"file", "listen_count", "duration", "artwork"})
		row.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number,
			track.File, track.ListenCount, track.Duration, track.Cover)
		return row
	}())

	result, err := GetTracks(1, db, "Pop")
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, track)
	}

	assert.Equal(t, tracks, result)
}

func TestGetAlbums(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	album := models.Album{
		Id:             1,
		Title:          "awa",
		Year:           2000,
		Artist:         "Земфира",
		ArtWork:        "awa",
		TracksCount:    1,
		TracksDuration: 1,
	}
	var albums []models.Album
	albums = append(albums, album)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
<<<<<<< HEAD
		SELECT a.id,
										       a.title,
										       a.year,
										       art.name,
										       a.artwork,
										       a.track_count,
										       SUM(t.duration) as tracksDuration
										FROM albums a
										         LEFT JOIN artists art ON art.id = a.artist
										         JOIN tracks t on t.album = a.id
										WHERE art.name = 'Земфира'
										GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
										LIMIT 1
=======
		SELECT DISTINCT ON (u.title) *
												FROM (SELECT DISTINCT ON (art.name) a.id,
												                                    a.title,
												                                    a.year,
												                                    art.name,
												                                    a.artwork,
												                                    a.track_count,
												                                    SUM(t.duration) as tracksDuration
												      FROM albums a
												               LEFT JOIN artists art ON art.id = a.artist
												               JOIN tracks t on t.album = a.id
												      WHERE RANDOM() <= 0.1
												      GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count
												      LIMIT 1) as u
>>>>>>> for_potikat
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"id", "title", "year", "artist", "artwork", "track_count", "tracksDuration"})
		row.AddRow(album.Id, album.Title, album.Year, album.Artist, album.ArtWork, album.TracksCount, album.TracksDuration)
		return row
	}())

	result, err := GetAlbums(1, db)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, album)
	}

	assert.Equal(t, albums, result)
}

func TestGetArtists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	artist := models.Artist{
		Id:     1,
		Name:   "awa",
		Avatar: "awa",
	}
	var artists []models.Artist
	artists = append(artists, artist)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
		SELECT artists.id, artists.name, artists.avatar FROM artists WHERE artists.avatar != 'frank_sinatra.jpg' LIMIT 1
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"id", "name", "avatar"})
		row.AddRow(artist.Id, artist.Name, artist.Avatar)
		return row
	}())

	result, err := GetArtists(1, db)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, artist)
	}

	assert.Equal(t, artists, result)
}

func TestGetPlaylists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	playlist := models.Playlist{
		Id:   1,
		Name: "awa",
		User: 1,
	}
	var playlists []models.Playlist
	playlists = append(playlists, playlist)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
		SELECT playlists.id, playlists.title, playlists.user FROM playlists LIMIT 1
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"id", "name", "user"})
		row.AddRow(playlist.Id, playlist.Name, playlist.User)
		return row
	}())

	result, err := GetPlaylists(1, db)
	if err != nil {
		t.Errorf("Error %s\n occurred during test case with %+v\n", err, playlist)
	}

	assert.Equal(t, playlists, result)
}

package utils

import (
	"2021_2_LostPointer/models"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

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
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"id", "title", "artist", "album", "explicit", "genre", "number",
			"file", "listen_count", "duration", "artwork"})
		row.AddRow(track.Id, track.Title, track.Artist, track.Album, track.Explicit, track.Genre, track.Number,
			track.File, track.ListenCount, track.Duration, track.Cover)
		return row
	}())

	result, err := GetTracks(1, db)
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
		Artist:         "awa",
		ArtWork:        "awa",
		TracksCount:    1,
		TracksDuration: 1,
	}
	var albums []models.Album
	albums = append(albums, album)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
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
		Id:   1,
		Name: "awa",
		Bio:  "awa",
		Avatar: "awa",
	}
	var artists []models.Artist
	artists = append(artists, artist)

	mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
		SELECT artists.id, artists.name, artists.bio, artists.avatar FROM artists LIMIT 1
	`))).WillReturnRows(func() *sqlmock.Rows {
		row := sqlmock.NewRows([]string{"id", "name", "bio", "avatar"})
		row.AddRow(artist.Id, artist.Name, artist.Bio, artist.Avatar)
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

package models

import (
	"database/sql"
	"fmt"
)

const tracksSelectionLimit = 50
const playlistsSelectionLimit = 4
const artistsSelectionLimit = 4

type Track struct {
	Id          int64
	Title       string
	Artist      string
	Album       string
	Explicit    bool
	Genre       string
	Number      int32
	File        string
	ListenCount uint64
}

type Artist struct {
	Id   int64
	Name string
	Bio  string
}

type Playlist struct {
	Id   int64
	Name string
	User int64
}

type SelectionFroHomePage struct {
	Tracks    []Track
	Playlists []Playlist
	Artists   []Artist
}

func GetTracks(amount int, db *sql.DB) ([]Track, error) {
	tracks := make([]Track, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT DISTINCT ON(alb.title, g.name) t.id, t.title, art.name, alb.title,  t.explicit, g.name, t.number, t.file, t.listen_count FROM tracks t
    											LEFT JOIN albums alb ON alb.id = t.album
    											LEFT JOIN artists art ON art.id = t.artist
    											LEFT JOIN genres g ON g.id = t.genre
    											LIMIT %d`, amount))
	if err != nil {
		return nil, err
	}

	var track Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Artist, &track.Album, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func GetPlaylists(amount int, db *sql.DB) ([]Playlist, error) {
	playlists := make([]Playlist, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM playlists LIMIT %d`, amount))
	if err != nil {
		return nil, err
	}

	var playlist Playlist
	for rows.Next() {
		if err := rows.Scan(&playlist.Id, &playlist.Name, &playlist.User); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func GetArtists(amount int, db *sql.DB) ([]Artist, error) {
	artists := make([]Artist, 0)

	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM artists LIMIT %d`, amount))
	if err != nil {
		return nil, err
	}

	var artist Artist
	for rows.Next() {
		if err := rows.Scan(&artist.Id, &artist.Name, &artist.Bio); err != nil {
			if artist.Bio == "" && artist.Name == "" {
				return nil, err
			}
		}
		artists = append(artists, artist)
	}

	return artists, nil
}

func GetSelectionForHomePage(db *sql.DB) (*SelectionFroHomePage, error) {
	var selectionForHomePage = new(SelectionFroHomePage)
	var err error

	if selectionForHomePage.Tracks, err = GetTracks(tracksSelectionLimit, db); err != nil {
		return nil, err
	}
	if selectionForHomePage.Playlists, err = GetPlaylists(playlistsSelectionLimit, db); err != nil {
		return nil, err
	}
	if selectionForHomePage.Artists, err = GetArtists(artistsSelectionLimit, db); err != nil {
		return nil, err
	}

	return selectionForHomePage, nil
}

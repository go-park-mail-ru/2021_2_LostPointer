package models

import (
	"database/sql"
	"fmt"
)

const tracksSelectionLimit = 50
const playlistsSelectionLimit = 4
const artistsSelectionLimit = 4

type Track struct {
	Id       int64
	Name     string
	Artist   int64
	Album    int64
	Explicit bool
	Genre    int64
	Number   int64
	File     string
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
	fmt.Println("1")
	rows, err := db.Query(fmt.Sprintf(`SELECT * FROM tracks LIMIT %d`, amount))
	if err != nil {
		fmt.Println("2")
		return nil, err
	}

	var track Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Name, &track.Artist, &track.Album,
			&track.Explicit, &track.Genre, &track.Number, &track.File); err != nil {
			fmt.Println(err)
			return nil, err
		}
		tracks = append(tracks, track)
	}
	fmt.Println(tracks)

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

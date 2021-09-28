package models

import (
	"database/sql"
	"fmt"
)

const tracksSelectionLimit = 10
const albumsSelectionLimit = 4
const playlistsSelectionLimit = 4
const artistsSelectionLimit = 4

type Track struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Artist      string `json:"artist"`
	Album       string `json:"album"`
	Explicit    bool   `json:"explicit"`
	Genre       string `json:"genre"`
	Number      int64  `json:"number"`
	File        string `json:"file"`
	ListenCount int64  `json:"listenCount"`
	Duration    int64  `json:"duration"`
	Cover       string `json:"cover"`
}

type Album struct {
	Id             int64  `json:"id"`
	Title          string `json:"title"`
	Year           int64  `json:"year"`
	Artist         string `json:"artist"`
	ArtWork        string `json:"artWork"`
	TrackCount     int64  `json:"trackCount"`
	TracksDuration int64  `json:"tracks_duration"`
}

type Artist struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type Playlist struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	User int64  `json:"user"`
}

type SelectionFroHomePage struct {
	Tracks    []Track    `json:"tracks"`
	Albums    []Album    `json:"albums"`
	Playlists []Playlist `json:"playlists"`
	Artists   []Artist   `json:"artists"`
}

func GetTracks(amount int, db *sql.DB) ([]Track, error) {
	tracks := make([]Track, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT DISTINCT ON(alb.title, g.name) t.id, t.title, art.name, alb.title,
												t.explicit, g.name, t.number, t.file, t.listen_count, t.duration, alb.artwork FROM tracks t
    											LEFT JOIN albums alb ON alb.id = t.album
    											LEFT JOIN artists art ON art.id = t.artist
    											LEFT JOIN genres g ON g.id = t.genre
    											LIMIT %d`, amount))
	if err != nil {
		fmt.Println(tracks)
		return nil, err
	}

	var track Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Artist, &track.Album, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Cover); err != nil {
			fmt.Println(err)

			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func GetAlbums(amount int, db *sql.DB) ([]Album, error) {
	albums := make([]Album, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT a.id, a.title, a.year, art.name, a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a
    											LEFT JOIN artists art ON art.id = a.artist
    											JOIN tracks t on t.album = a.id
    											GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count LIMIT %d`, amount))
	if err != nil {
		fmt.Println(err)

		return nil, err
	}

	var album Album
	for rows.Next() {
		if err := rows.Scan(&album.Id, &album.Title, &album.Year, &album.Artist, &album.ArtWork, &album.TrackCount, &album.TracksDuration); err != nil {
			fmt.Println(err)
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
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
	if selectionForHomePage.Albums, err = GetAlbums(albumsSelectionLimit, db); err != nil {
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

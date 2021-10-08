package utils

import (
	"2021_2_LostPointer/models"
	"database/sql"
	"fmt"
)

const TracksSelectionLimit = 10
const AlbumsSelectionLimit = 4
const PlaylistsSelectionLimit = 4
const ArtistsSelectionLimit = 4

func GetTracks(amount int, db *sql.DB) ([]models.Track, error) {
	tracks := make([]models.Track, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT DISTINCT ON (alb.title) t.id,
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
												LIMIT %d`, amount))
	if err != nil {
		return nil, err
	}

	var track models.Track

	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Artist, &track.Album, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Cover); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func GetAlbums(amount int, db *sql.DB) ([]models.Album, error) {
	albums := make([]models.Album, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT DISTINCT ON (u.title) *
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
												      LIMIT %d) as u`, amount))

	if err != nil {
		return nil, err
	}

	var album models.Album
	for rows.Next() {
		if err := rows.Scan(&album.Id, &album.Title, &album.Year, &album.Artist, &album.ArtWork, &album.TracksCount, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func GetArtists(amount int, db *sql.DB) ([]models.Artist, error) {
	artists := make([]models.Artist, 0)

	rows, err := db.Query(fmt.Sprintf(`SELECT artists.id, artists.name, artists.bio, artists.avatar FROM artists LIMIT %d`, amount))
	if err != nil {
		return nil, err
	}

	var artist models.Artist
	for rows.Next() {
		if err := rows.Scan(&artist.Id, &artist.Name, &artist.Bio, &artist.Avatar); err != nil {
			if artist.Bio == "" && artist.Name == "" {
				return nil, err
			}
		}
		fmt.Println(artist)
		artists = append(artists, artist)
	}

	return artists, nil
}

func GetPlaylists(amount int, db *sql.DB) ([]models.Playlist, error) {
	playlists := make([]models.Playlist, 0)
	rows, err := db.Query(fmt.Sprintf(`SELECT playlists.id, playlists.title, playlists.user FROM playlists LIMIT %d`, amount))
	if err != nil {
		return nil, err
	}

	var playlist models.Playlist
	for rows.Next() {
		if err := rows.Scan(&playlist.Id, &playlist.Name, &playlist.User); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func GetSelectionForHomePage(db *sql.DB) (*models.SelectionFroHomePage, error) {
	var selectionForHomePage = new(models.SelectionFroHomePage)
	var err error

	if selectionForHomePage.Tracks, err = GetTracks(TracksSelectionLimit, db); err != nil {
		return nil, err
	}
	if selectionForHomePage.Albums, err = GetAlbums(AlbumsSelectionLimit, db); err != nil {
		return nil, err
	}
	if selectionForHomePage.Playlists, err = GetPlaylists(PlaylistsSelectionLimit, db); err != nil {
		return nil, err
	}
	if selectionForHomePage.Artists, err = GetArtists(ArtistsSelectionLimit, db); err != nil {
		return nil, err
	}

	return selectionForHomePage, nil
}

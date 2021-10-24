package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql"
)

type MusicRepository struct {
	Database *sql.DB
}

func NewMusicRepository(db *sql.DB) MusicRepository {
	return MusicRepository{Database: db}
}

func (musicRepository MusicRepository) GetRandomTracks(amount int, isAuthorized bool) ([]models.Track, error) {
	rows, err := musicRepository.Database.Query("SELECT tracks.id, tracks.title, art.name, alb.title, explicit, "+
		"g.name, number, file, listen_count, duration, lossless, alb.artwork as cover FROM tracks "+
		"LEFT JOIN genres g ON tracks.genre = g.id "+
		"LEFT JOIN albums alb ON tracks.album = alb.id "+
		"LEFT JOIN artists art ON tracks.artist = art.id ORDER BY RANDOM() LIMIT $1", amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	tracks := make([]models.Track, 0, 10)
	var track models.Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Artist, &track.Album, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Lossless, &track.Cover); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (musicRepository MusicRepository) GetRandomAlbums(amount int) ([]models.Album, error) {
	rows, err := musicRepository.Database.Query("SELECT a.id, a.title, a.year, art.name, "+
		"a.artwork, a.track_count, SUM(t.duration) as tracksDuration FROM albums a "+
		"LEFT JOIN artists art ON art.id = a.artist "+
		"JOIN tracks t on t.album = a.id "+
		"WHERE art.name NOT LIKE '%Frank Sinatra%' "+
		"GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count "+
		"ORDER BY RANDOM() "+
		"LIMIT $1", amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	albums := make([]models.Album, 0, 10)
	var album models.Album
	for rows.Next() {
		if err := rows.Scan(&album.Id, &album.Title, &album.Year, &album.Artist, &album.Artwork, &album.TracksCount, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func (musicRepository MusicRepository) GetRandomArtists(amount int) ([]models.Artist, error) {
	rows, err := musicRepository.Database.Query("SELECT artists.id, artists.name, artists.avatar FROM artists "+
		"WHERE artists.name NOT LIKE '%Frank Sinatra%' ORDER BY RANDOM() LIMIT $1", amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	artists := make([]models.Artist, 0, 10)
	var artist models.Artist
	for rows.Next() {
		if err := rows.Scan(&artist.Id, &artist.Name, &artist.Avatar); err != nil {
			return nil, err
		}
		artists = append(artists, artist)
	}
	return artists, nil
}

func (musicRepository MusicRepository) GetRandomPlaylists(amount int) ([]models.Playlist, error) {
	rows, err := musicRepository.Database.Query("SELECT playlists.id, playlists.title, playlists.user "+
		"FROM playlists LIMIT $1", amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	playlists := make([]models.Playlist, 0, 10)
	var playlist models.Playlist
	for rows.Next() {
		if err := rows.Scan(&playlist.Id, &playlist.Name, &playlist.User); err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql"
)

type ArtistRepository struct {
	Database *sql.DB
}

func NewArtistRepository(db *sql.DB) ArtistRepository {
	return ArtistRepository{Database: db}
}

func (artistRepository ArtistRepository) Get(id int) (*models.Artist, error) {
	artist := new(models.Artist)
	err := artistRepository.Database.QueryRow("SELECT art.id, art.name, art.avatar FROM artists art "+
		"WHERE art.id = $1 "+
		"GROUP BY art.id", id).Scan(&artist.Id, &artist.Name, &artist.Avatar)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (artistRepository ArtistRepository) GetTracks(id int, isAuthorized bool, amount int) ([]models.Track, error) {
	var tracks []models.Track
	trackRows, err := artistRepository.Database.Query("SELECT tracks.id, tracks.title, explicit, "+
		"g.name, number, file, listen_count, duration, lossless, alb.id, alb.title, alb.artwork as cover FROM tracks "+
		"LEFT JOIN genres g ON tracks.genre = g.id "+
		"LEFT JOIN albums alb ON tracks.album = alb.id "+
		"LEFT JOIN artists art ON tracks.artist = art.id "+
		"WHERE tracks.artist = $1 "+
		"ORDER BY tracks.listen_count LIMIT $2", id, amount)
	if err != nil {
		return nil, err
	}
	defer func(trackRows *sql.Rows) {
		_ = trackRows.Close()
	}(trackRows)

	var track models.Track
	for trackRows.Next() {
		if err := trackRows.Scan(&track.Id, &track.Title, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Lossless, &track.Album.Id,
			&track.Album.Title, &track.Album.Artwork); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (artistRepository ArtistRepository) GetAlbums(id int, amount int) ([]models.Album, error) {
	var albums []models.Album
	albumRows, err := artistRepository.Database.Query("SELECT a.id, a.title, a.artwork, a.year, SUM(t.duration) AS "+
		"tracksDuration FROM albums a "+
		"JOIN tracks t on t.album = a.id "+
		"WHERE a.artist = $1 "+
		"GROUP BY a.id, a.title, a.artwork, a.year "+
		"ORDER BY a.year DESC LIMIT $2", id, amount)
	if err != nil {
		return nil, err
	}
	defer func(albumRows *sql.Rows) {
		_ = albumRows.Close()
	}(albumRows)

	var album models.Album
	for albumRows.Next() {
		if err := albumRows.Scan(&album.Id, &album.Title, &album.Artwork, &album.Year, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}
func (artistRepository ArtistRepository) GetRandom(amount int) ([]models.Artist, error) {
	rows, err := artistRepository.Database.Query("SELECT artists.id, artists.name, artists.avatar FROM artists "+
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

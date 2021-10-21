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
	trackRows, err := artistRepository.Database.Query("SELECT t.id, t.title, explicit, file, duration, lossless, "+
		"alb.artwork FROM tracks t "+
		"LEFT JOIN albums alb ON alb.id = t.album "+
		"WHERE t.artist = $1 "+
		"ORDER BY t.listen_count LIMIT $2", id, amount)
	if err != nil {
		return nil, err
	}
	defer func(trackRows *sql.Rows) {
		_ = trackRows.Close()
	}(trackRows)

	var track models.Track
	for trackRows.Next() {
		if err := trackRows.Scan(&track.Id, &track.Title, &track.Explicit, &track.File, &track.Duration, &track.Lossless, &track.Cover); err != nil {
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
		if err := albumRows.Scan(&album.Id, &album.Title, &album.ArtWork, &album.Year, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

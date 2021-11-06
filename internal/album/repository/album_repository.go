package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql"
)

type AlbumRepository struct {
	Database *sql.DB
}

func NewAlbumRepository(db *sql.DB) AlbumRepository {
	return AlbumRepository{Database: db}
}

func (albumRepository *AlbumRepository) GetRandom(amount int) ([]models.Album, error) {
	rows, err := albumRepository.Database.Query("SELECT a.id, a.title, a.year, art.name, "+
		"a.artwork, a.track_count, SUM(t.duration) as tracksDuration, a.artwork_color FROM albums a "+
		"LEFT JOIN artists art ON art.id = a.artist "+
		"JOIN tracks t on t.album = a.id "+
		"GROUP BY a.id, a.title, a.year, art.name, a.artwork, a.track_count "+
		"ORDER BY RANDOM() "+
		"LIMIT $1", amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	albums := make([]models.Album, 0, amount)
	for rows.Next() {
		var album models.Album
		if err := rows.Scan(&album.Id, &album.Title, &album.Year, &album.Artist, &album.Artwork, &album.TracksCount, &album.TracksDuration, &album.ArtWorkColor); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func (albumRepository *AlbumRepository) GetByArtistID(id, amount int) ([]models.Album, error) {
	var albums []models.Album
	albumRows, err := albumRepository.Database.Query("SELECT a.id, a.title, a.artwork, a.year, SUM(t.duration) AS "+
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

	for albumRows.Next() {
		var album models.Album
		if err := albumRows.Scan(&album.Id, &album.Title, &album.Artwork, &album.Year, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

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

func (albumRepository AlbumRepository) GetRandom(amount int) ([]models.Album, error) {
	rows, err := albumRepository.Database.Query("SELECT a.id, a.title, a.year, art.name, "+
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
		if err := rows.Scan(&album.Id, &album.Title, &album.Year, &album.Artist, &album.ArtWork, &album.TracksCount, &album.TracksDuration); err != nil {
			return nil, err
		}
		albums = append(albums, album)
	}

	return albums, nil
}

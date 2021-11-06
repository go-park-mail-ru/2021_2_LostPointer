package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql"
	"log"
)

type PlaylistRepository struct {
	Database *sql.DB
}

func NewPlaylistRepository(db *sql.DB) PlaylistRepository {
	return PlaylistRepository{Database: db}
}

func (playlistRepository PlaylistRepository) Get(amount int, id int) ([]models.Playlist, error) {
	rows, err := playlistRepository.Database.Query("SELECT playlists.id, playlists.title, playlists.user "+
		"FROM playlists WHERE playlists.user = $1 LIMIT $2", id, amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			log.Println(err)
		}
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

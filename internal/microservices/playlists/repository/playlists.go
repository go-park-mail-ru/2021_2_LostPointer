package repository

import (
	"database/sql"
	"log"

	"github.com/kennygrant/sanitize"

	"2021_2_LostPointer/internal/microservices/playlists/proto"
)

type PlaylistsStorage struct {
	db *sql.DB
}

func NewPlaylistsStorage(db *sql.DB) *PlaylistsStorage {
	return &PlaylistsStorage{db: db}
}

func (storage *PlaylistsStorage) CreatePlaylist(userID int64, title string) (*proto.CreatePlaylistResponse, error) {
	query := `INSERT INTO playlists(title, user_id) VALUES ($1, $2) RETURNING id`

	var id int64
	title = sanitize.HTML(title)
	err := storage.db.QueryRow(query, title, userID).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &proto.CreatePlaylistResponse{PlaylistID: id}, nil
}

func (storage *PlaylistsStorage) UpdatePlaylist(playlistID int64, title string) error {
	query := `UPDATE playlists SET title=$1 WHERE id=$2`

	title = sanitize.HTML(title)
	err := storage.db.QueryRow(query, title, playlistID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage *PlaylistsStorage) DeletePlaylist(playlistID int64) error {
	query := `DELETE FROM playlists WHERE id=$1`

	err := storage.db.QueryRow(query, playlistID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage *PlaylistsStorage) AddTrack(playlistID int64, trackID int64) error {
	query := `INSERT INTO playlist_tracks(playlist, track) VALUES ($1, $2)`

	err := storage.db.QueryRow(query, playlistID, trackID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage *PlaylistsStorage) IsAdded(playlistID int64, trackID int64) (bool, error) {
	query := `SELECT * FROM playlist_tracks WHERE playlist=$1 AND track=$2`

	rows, err := storage.db.Query(query, playlistID, trackID)
	if err != nil {
		return true, err
	}
	err = rows.Err()
	if err != nil {
		return true, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (storage *PlaylistsStorage) IsOwner(playlistID int64, userID int64) (bool, error) {
	query := `SELECT * FROM playlists WHERE id=$1 AND user_id=$2`

	rows, err := storage.db.Query(query, playlistID, userID)
	if err != nil {
		return true, err
	}
	err = rows.Err()
	if err != nil {
		return true, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (storage *PlaylistsStorage) UserPlaylists(userID int64) ([]*proto.Playlist, error) {


	return nil, nil
}

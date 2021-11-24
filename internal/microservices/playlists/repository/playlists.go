package repository

import (
	"2021_2_LostPointer/internal/constants"
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

func (storage *PlaylistsStorage) CreatePlaylist(userID int64, title string, artwork string, artworkColor string, isPublic bool) (*proto.CreatePlaylistResponse, error) {
	query := `INSERT INTO playlists(title, user_id, artwork, artwork_color, is_public) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id int64
	title = sanitize.HTML(title)
	err := storage.db.QueryRow(query, title, userID, artwork, artworkColor, isPublic).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &proto.CreatePlaylistResponse{PlaylistID: id}, nil
}

func (storage *PlaylistsStorage) GetOldPlaylistSettings(playlistID int64) (string, error) {
	query := `SELECT artwork FROM playlists WHERE id=$1`

	var artwork string
	err := storage.db.QueryRow(query, playlistID).Scan(&artwork)
	if err != nil {
		return "", err
	}

	return artwork, nil
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

//nolint:rowserrcheck
func (storage *PlaylistsStorage) IsAdded(playlistID int64, trackID int64) (bool, error) {
	query := `SELECT * FROM playlist_tracks WHERE playlist=$1 AND track=$2`

	rows, err := storage.db.Query(query, playlistID, trackID)
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

func (storage *PlaylistsStorage) DeleteTrack(playlistID int64, trackID int64) error {
	query := `DELETE FROM playlist_tracks WHERE playlist=$1 AND track=$2`

	err := storage.db.QueryRow(query, playlistID, trackID).Err()
	if err != nil {
		return err
	}

	return nil
}

//nolint:rowserrcheck
func (storage *PlaylistsStorage) IsOwner(playlistID int64, userID int64) (bool, error) {
	query := `SELECT * FROM playlists WHERE id=$1 AND (user_id=$2 OR is_public=true)`

	rows, err := storage.db.Query(query, playlistID, userID)
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

func (storage *PlaylistsStorage) DoesPlaylistExist(playlistID int64) (bool, error) {
	query := `SELECT * FROM playlists WHERE id=$1`

	rows, err := storage.db.Query(query, playlistID)
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

func (storage *PlaylistsStorage) UpdatePlaylistTitle(playlistID int64, title string) error {
	query := `UPDATE playlists SET title=$1 WHERE id=$2`

	title = sanitize.HTML(title)
	err := storage.db.QueryRow(query, title, playlistID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage *PlaylistsStorage) UpdatePlaylistArtwork(playlistID int64, artwork string, artworkColor string) error {
	query := `UPDATE playlists SET artwork=$1, artwork_color=$2 WHERE id=$3`

	err := storage.db.QueryRow(query, artwork, artworkColor, playlistID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage *PlaylistsStorage) DeletePlaylistArtwork(playlistID int64) error {
	query := `UPDATE playlists SET artwork=$1, artwork_color=$2 WHERE id=$3`

	err := storage.db.QueryRow(query, constants.PlaylistArtworkDefaultFilename, constants.PlaylistArtworkDefaultColor, playlistID).Err()
	if err != nil {
		return err
	}

	return nil
}

func (storage *PlaylistsStorage) UpdatePlaylistAccess(playlistID int64, isPublic bool) error {
	query := `UPDATE playlists SET is_public=$1 WHERE id=$2`

	err := storage.db.QueryRow(query, isPublic, playlistID).Err()
	if err != nil {
		return err
	}

	return nil
}

package repository

import (
	"2021_2_LostPointer/internal/models"
	"database/sql"
)

type TrackRepository struct {
	Database *sql.DB
}

func NewTrackRepository(db *sql.DB) TrackRepository {
	return TrackRepository{Database: db}
}

func (trackRepository *TrackRepository) GetRandom(amount int, isAuthorized bool) ([]models.Track, error) {
	rows, err := trackRepository.Database.Query("SELECT tracks.id, tracks.title, explicit, "+
		"g.name, number, file, listen_count, duration, lossless, alb.id, alb.title, alb.artwork, art.id, art.name FROM tracks "+
		"LEFT JOIN genres g ON tracks.genre = g.id "+
		"LEFT JOIN albums alb ON tracks.album = alb.id "+
		"LEFT JOIN artists art ON tracks.artist = art.id ORDER BY RANDOM() LIMIT $1", amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	tracks := make([]models.Track, 0, amount)
	for rows.Next() {
		var track models.Track
		track.Album = new(models.Album)
		track.Artist = new(models.Artist)
		if err := rows.Scan(&track.Id, &track.Title, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Lossless, &track.Album.Id,
			&track.Album.Title, &track.Album.Artwork, &track.Artist.Id, &track.Artist.Name); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (trackRepository *TrackRepository) IncrementListenCount(id int64) error {
	err := trackRepository.Database.QueryRow(`UPDATE tracks SET listen_count = listen_count + 1 WHERE id=$1`, id).Err()
	if err != nil {
		return err
	}

	return nil
}

func (trackRepository *TrackRepository) GetByAlbumID(albumID int, isAuthorized bool) ([]models.Track, error) {
	rows, err := trackRepository.Database.Query("SELECT tracks.id, tracks.title, explicit, "+
		"g.name, number, file, listen_count, duration, lossless, art.id, art.name FROM tracks "+
		"LEFT JOIN genres g ON tracks.genre = g.id "+
		"LEFT JOIN albums alb ON tracks.album = alb.id "+
		"LEFT JOIN artists art ON tracks.artist = art.id WHERE alb.id = $1 ORDER BY tracks.listen_count DESC", albumID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	tracks := make([]models.Track, 0, 10)
	for rows.Next() {
		var track models.Track
		track.Artist = new(models.Artist)
		if err := rows.Scan(&track.Id, &track.Title, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Lossless,
			&track.Artist.Id, &track.Artist.Name); err != nil {
			return nil, err
		}
		if !isAuthorized {
			track.File = ""
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (trackRepository *TrackRepository) GetByArtistID(artistID, amount int, isAuthorized bool) ([]models.Track, error) {
	rows, err := trackRepository.Database.Query("SELECT tracks.id, tracks.title, explicit, "+
		"g.name, number, file, listen_count, duration, lossless, alb.id, alb.title, alb.artwork FROM tracks "+
		"LEFT JOIN genres g ON tracks.genre = g.id "+
		"LEFT JOIN albums alb ON tracks.album = alb.id "+
		"LEFT JOIN artists art ON tracks.artist = art.id WHERE art.id = $1 ORDER BY tracks.listen_count DESC "+
		"LIMIT $2", artistID, amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	tracks := make([]models.Track, 0, amount)
	for rows.Next() {
		var track models.Track
		track.Album = new(models.Album)
		if err := rows.Scan(&track.Id, &track.Title, &track.Explicit, &track.Genre,
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

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

func (trackRepository TrackRepository) GetRandom(amount int, isAuthorized bool) ([]models.Track, error) {
	rows, err := trackRepository.Database.Query("SELECT tracks.id, tracks.title, art.name, alb.title, explicit, "+
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
package repository

import (
	"2021_2_LostPointer/internal/microservices/music/proto"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/pkg/wrapper"
	"database/sql"
	"log"
)

type MusicStorage struct {
	db *sql.DB
}

func NewMusicStorage(db *sql.DB) *MusicStorage {
	return &MusicStorage{db: db}
}

func (storage *MusicStorage) RandomTracks(amount int64, isAuthorized bool) (*proto.Tracks, error) {
	query := `SELECT ` +
		wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t") + ", " +
		wrapper.Wrapper([]string{"id", "title", "artwork"}, "alb") + ", " +
		wrapper.Wrapper([]string{"id", "name"}, "art") + ", " +
		wrapper.Wrapper([]string{"name"}, "g") +
		`
		FROM tracks t
		LEFT JOIN genres g ON t.genre = g.id
		LEFT JOIN albums alb ON t.album = alb.id
		LEFT JOIN artists art ON t.artist = art.id
		ORDER BY RANDOM() DESC LIMIT $1`

	rows, err := storage.db.Query(query, amount)
	if err != nil {
		return nil, err
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal("Error occurred during closing rows")
		}
	}()

	// tracks := make([]*proto.Track, 0, constants.TracksSearchAmount)
	
	for rows.Next() {
		//var track proto.Track
		var track models.Track
		if err = rows.Scan(&track.ID, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.ID, &track.Album.Title, &track.Album.Artwork,
			&track.Artist.ID, &track.Artist.Name, &track.Genre); err != nil {
			return nil, err
		}
		//if !isAuthorized {
		//	track.File = ""
		//}
		log.Println(track)
		// tracks = append(tracks, track)
	}

	//tracksList := &proto.Tracks{Tracks: tracks}

	return nil, nil
}
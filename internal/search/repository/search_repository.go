package repository

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/pkg/wrapper"
	"database/sql"
	"log"
	"os"
)

type MusicInfoStorage struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) MusicInfoStorage {
	return MusicInfoStorage{db: db}
}


func (storage *MusicInfoStorage) TracksByFullWord(text string) ([]models.Track, error) {
	log.Println(wrapper.Wrapper([]string{"id", "title", "explicit", "number", "file", "listen_count", "duration", "lossless"}, "t"))
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
		WHERE t.id IN (
		    SELECT track
		    FROM search
		    WHERE concatenation @@ plainto_tsquery($1)
		)
		ORDER BY listen_count DESC LIMIT $2`

	log.Println(query)

	rows, err := storage.db.Query(query, text, constants.TracksSearchAmount)
	if err != nil {
		log.Println("FIRST")
		return nil, err
	}

	tracks := make([]models.Track, 0, constants.TracksSearchAmount)
	var track models.Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.Id, &track.Album.Title, &track.Album.Artwork,
			&track.Artist.Id, &track.Artist.Name, &track.Genre); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (storage *MusicInfoStorage) TracksByPartial(text string) ([]models.Track, error) {
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
		WHERE t.id IN (
		    SELECT track
		    FROM test
		    WHERE concat ILIKE $1
		)
		ORDER BY listen_count DESC LIMIT $2`

	rows, err := storage.db.Query(query, "%"+text+"%", constants.TracksSearchAmount)
	if err != nil {
		log.Println("SECOND")
		return nil, err
	}

	tracks := make([]models.Track, 0, constants.TracksSearchAmount)
	var track models.Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Explicit, &track.Number, &track.File, &track.ListenCount,
			&track.Duration, &track.Lossless, &track.Album.Id, &track.Album.Title, &track.Album.Artwork,
			&track.Artist.Id, &track.Artist.Name, &track.Genre); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (storage *MusicInfoStorage) Artists(text string) ([]models.ArtistShort, error) {
	rows, err := storage.db.Query(`
		SELECT id, name, avatar
		FROM artists
		WHERE name ILIKE $1
	`, "%" + text + "%")
	if err != nil {
		return nil, err
	}

	artists := make([]models.ArtistShort, 0, constants.ArtistsSearchAmount)
	var artist models.ArtistShort
	for rows.Next() {
		if err := rows.Scan(&artist.ID, &artist.Name, &artist.Avatar); err != nil {
			return nil, err
		}
		artist.Avatar = os.Getenv("ARTISTS_ROOT_PREFIX") + artist.Avatar + constants.LittleAvatarPostfix
		artists = append(artists, artist)
	}

	return artists, nil
}

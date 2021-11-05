package repository

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"database/sql"
	"os"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) SearchRepository {
	return SearchRepository{db: db}
}

func (searchR SearchRepository) SearchRelevantTracksByFullWord(text string) ([]models.Track, error) {
	rows, err := searchR.db.Query(`
		SELECT t.id, t.title, explicit, g.name, number, file, listen_count, duration, lossless, alb.id, alb.title,
		       alb.artwork, art.id, art.name
		FROM tracks t
		LEFT JOIN genres g ON t.genre = g.id
		LEFT JOIN albums alb ON t.album = alb.id
		LEFT JOIN artists art ON t.artist = art.id
		WHERE t.id IN (
		    SELECT track
		    FROM search
		    WHERE concatenation @@ plainto_tsquery($1)
		)
		ORDER BY listen_count DESC LIMIT $2
	`, text, constants.TracksSearchAmount)
	if err != nil {
		return nil, err
	}

	tracks := make([]models.Track, 0, constants.TracksSearchAmount)
	var track models.Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Lossless, &track.Album.Id,
			&track.Album.Title, &track.Album.Artwork, &track.Artist.Id, &track.Artist.Name); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (searchR SearchRepository) SearchRelevantTracksByPartial(text string) ([]models.Track, error) {
	rows, err := searchR.db.Query(`
		SELECT t.id, t.title, explicit, g.name, number, file, listen_count, duration, lossless, alb.id, alb.title,
		       alb.artwork, art.id, art.name
		FROM tracks t
		LEFT JOIN genres g ON t.genre = g.id
		LEFT JOIN albums alb ON t.album = alb.id
		LEFT JOIN artists art ON t.artist = art.id
		WHERE t.id IN (
		    SELECT track
		    FROM test
		    WHERE concat ILIKE $1
		)
		ORDER BY listen_count DESC LIMIT $2
	`, "%"+text+"%", constants.TracksSearchAmount)
	if err != nil {
		return nil, err
	}

	tracks := make([]models.Track, 0, constants.TracksSearchAmount)
	var track models.Track
	for rows.Next() {
		if err := rows.Scan(&track.Id, &track.Title, &track.Explicit, &track.Genre,
			&track.Number, &track.File, &track.ListenCount, &track.Duration, &track.Lossless, &track.Album.Id,
			&track.Album.Title, &track.Album.Artwork, &track.Artist.Id, &track.Artist.Name); err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}

	return tracks, nil
}

func (searchR SearchRepository) SearchRelevantArtists(text string) ([]models.ArtistShort, error) {
	rows, err := searchR.db.Query(`
		SELECT id, name, avatar
		FROM artists
		WHERE name ILIKE $1
	`, "%"+text+"%")
	if err != nil {
		return nil, err
	}

	artists := make([]models.ArtistShort, 0, constants.ArtistsSearchAmount)
	var artist models.ArtistShort
	for rows.Next() {
		if err := rows.Scan(&artist.Id, &artist.Name, &artist.Avatar); err != nil {
			return nil, err
		}
		artist.Avatar = os.Getenv("ARTISTS_ROOT_PREFIX") + artist.Avatar + constants.LittleAvatarPostfix
		artists = append(artists, artist)
	}

	return artists, nil
}

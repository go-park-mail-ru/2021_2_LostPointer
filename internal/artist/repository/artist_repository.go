package repository

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"database/sql"
	"os"
)

type ArtistRepository struct {
	Database *sql.DB
}

func NewArtistRepository(db *sql.DB) ArtistRepository {
	return ArtistRepository{Database: db}
}

func (artistRepository *ArtistRepository) Get(id int) (*models.Artist, error) {
	artist := new(models.Artist)
	var video string
	err := artistRepository.Database.QueryRow("SELECT art.id, art.name, art.avatar, art.video FROM artists art "+
		"WHERE art.id = $1 "+
		"GROUP BY art.id", id).Scan(&artist.Id, &artist.Name, &artist.Avatar, &video)
	artist.Avatar = os.Getenv("ARTISTS_ROOT_PREFIX") + artist.Avatar + constants.BigArtistPostfix
	if len(video) != 1 {
		artist.Video = os.Getenv("MOV_ROOT_PREFIX") + video + constants.VideoPostfix
	}
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (artistRepository *ArtistRepository) GetRandom(amount int) ([]models.Artist, error) {
	rows, err := artistRepository.Database.Query("SELECT artists.id, artists.name, artists.avatar, artists.video FROM artists "+
		"ORDER BY RANDOM() LIMIT $1", amount)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	artists := make([]models.Artist, 0, amount)

	for rows.Next() {
		var artist models.Artist
		var video string
		if err := rows.Scan(&artist.Id, &artist.Name, &artist.Avatar, &video); err != nil {
			return nil, err
		}
		artist.Avatar = os.Getenv("ARTISTS_ROOT_PREFIX") + artist.Avatar + constants.BigArtistPostfix
		if len(video) > 1 {
			artist.Video = os.Getenv("MOV_ROOT_PREFIX") + video + constants.VideoPostfix
		}
		artists = append(artists, artist)
	}
	return artists, nil
}

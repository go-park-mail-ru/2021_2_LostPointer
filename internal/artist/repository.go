package artist

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/artist_repo_db_mock.go -pkg mock . ArtistRepository:MockArtistRepository
type ArtistRepository interface {
	Get(id int) (*models.Artist, error)
	GetTracks(id int, isAuthorized bool) ([]models.Track, error)
	GetAlbums(id int) ([]models.Album, error)
}

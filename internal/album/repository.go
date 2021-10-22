package album

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/album_repo_db_mock.go -pkg mock . AlbumRepository:MockAlbumRepository
type AlbumRepository interface {
	GetRandom(amount int) ([]models.Album, error)
}

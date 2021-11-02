package album

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/album_usecase_db_mock.go -pkg mock . AlbumUseCase:MockAlbumUseCase
type AlbumUseCase interface {
	GetHome(amount int) ([]models.Album, *models.CustomError)
}


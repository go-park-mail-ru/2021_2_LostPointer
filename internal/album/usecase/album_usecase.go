package usecase

import (
	"2021_2_LostPointer/internal/album"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/utils/constants"
	"net/http"
)


type AlbumUseCase struct {
	AlbumRepository album.AlbumRepository
}

func NewAlbumUseCase(albumRepository album.AlbumRepository) AlbumUseCase {
	return AlbumUseCase{AlbumRepository: albumRepository}
}

func (albumUseCase AlbumUseCase) GetHome(amount int) ([]models.Album, *models.CustomError) {
	albums, err := albumUseCase.AlbumRepository.GetRandom(amount)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}

	return albums, nil
}

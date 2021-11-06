package usecase

import (
	"2021_2_LostPointer/internal/artist"
	_ "2021_2_LostPointer/internal/artist"
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"net/http"
)

type ArtistUseCase struct {
	ArtistRepository artist.ArtistRepository
}

func NewArtistUseCase(artistRepository artist.ArtistRepository) ArtistUseCase {
	return ArtistUseCase{ArtistRepository: artistRepository}
}

func (artistUseCase *ArtistUseCase) GetProfile(id int) (*models.Artist, *models.CustomError) {
	var err error
	art, err := artistUseCase.ArtistRepository.Get(id)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}

	return art, nil
}

func (artistUseCase *ArtistUseCase) GetHome(amount int) ([]models.Artist, *models.CustomError) {
	artists, err := artistUseCase.ArtistRepository.GetRandom(amount)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}

	return artists, nil
}

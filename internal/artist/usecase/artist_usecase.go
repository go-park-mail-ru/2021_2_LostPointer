package usecase

import (
	"2021_2_LostPointer/internal/artist"
	_ "2021_2_LostPointer/internal/artist"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/utils/constants"
	"net/http"
)

type ArtistUseCase struct {
	ArtistRepository artist.ArtistRepository
}

func NewArtistUseCase(artistRepository artist.ArtistRepository) ArtistUseCase {
	return ArtistUseCase{ArtistRepository: artistRepository}
}

func (artistUseCase ArtistUseCase) GetProfile(id int, isAuthorized bool) (*models.Artist, *models.CustomError) {
	var err error
	art, err := artistUseCase.ArtistRepository.Get(id)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}
	art.Tracks, err = artistUseCase.ArtistRepository.GetTracks(id, isAuthorized, constants.TracksDefaultAmountForArtist)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}
	art.Albums, err = artistUseCase.ArtistRepository.GetAlbums(id, constants.AlbumsDefaultAmountForArtist)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}

	return art, nil
}

func (artistUseCase ArtistUseCase) GetHome(amount int) ([]models.Artist, *models.CustomError) {
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

package usecase

import (
	"2021_2_LostPointer/internal/artist"
	_ "2021_2_LostPointer/internal/artist"
	"2021_2_LostPointer/internal/models"
	"net/http"
)

const DatabaseNotResponding = "Database not responding"

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
			Message:       DatabaseNotResponding,
		}
	}
	art.Tracks, err = artistUseCase.ArtistRepository.GetTracks(id, isAuthorized)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       DatabaseNotResponding,
		}
	}
	art.Albums, err = artistUseCase.ArtistRepository.GetAlbums(id)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       DatabaseNotResponding,
		}
	}

	return art, nil
}

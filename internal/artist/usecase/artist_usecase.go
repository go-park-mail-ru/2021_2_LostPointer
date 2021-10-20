package usecase

import (
	"2021_2_LostPointer/internal/artist/repository"
	"2021_2_LostPointer/internal/models"
	"net/http"
)

const DatabaseNotResponding = "Database not responding"

type ArtistUseCase struct {
	ArtistRepository repository.ArtistRepository
}

func NewArtistUseCase(artistRepository repository.ArtistRepository) ArtistUseCase {
	return ArtistUseCase{ArtistRepository: artistRepository}
}

func (artistUseCase ArtistUseCase) GetProfile(id int, isAuthorized bool) (*models.Artist, *models.CustomError) {
	var err error
	artist, err := artistUseCase.ArtistRepository.Get(id)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       DatabaseNotResponding,
		}
	}
	artist.Tracks, err = artistUseCase.ArtistRepository.GetTracks(id, isAuthorized)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       DatabaseNotResponding,
		}
	}
	artist.Albums, err = artistUseCase.ArtistRepository.GetAlbums(id)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       DatabaseNotResponding,
		}
	}

	return artist, nil
}

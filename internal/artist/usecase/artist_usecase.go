package usecase

import (
	"2021_2_LostPointer/internal/artist"
	_ "2021_2_LostPointer/internal/artist"
	"2021_2_LostPointer/internal/models"
	"net/http"
)

const DatabaseNotResponding = "Database not responding"
const TracksDefaultAmount = 20
const AlbumsDefaultAmount = 8


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
	art.Tracks, err = artistUseCase.ArtistRepository.GetTracks(id, isAuthorized, TracksDefaultAmount)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       DatabaseNotResponding,
		}
	}
	art.Albums, err = artistUseCase.ArtistRepository.GetAlbums(id, AlbumsDefaultAmount)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       DatabaseNotResponding,
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
			Message:       DatabaseNotResponding,
		}
	}

	return artists, nil
}

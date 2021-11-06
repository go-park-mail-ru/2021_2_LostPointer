package usecase

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/track"
	"fmt"
	"net/http"
)

type TrackUseCase struct {
	TrackRepository track.TrackRepository
}

func NewTrackUseCase(trackRepository track.TrackRepository) TrackUseCase {
	return TrackUseCase{TrackRepository: trackRepository}
}

func (trackUseCase *TrackUseCase) GetHome(amount int, isAuthorized bool) ([]models.Track, *models.CustomError) {
	tracks, err := trackUseCase.TrackRepository.GetRandom(amount, isAuthorized)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}

	return tracks, nil
}

func (trackUseCase *TrackUseCase) IncrementListenCount(id int64) *models.CustomError {
	err := trackUseCase.TrackRepository.IncrementListenCount(id)
	if err != nil {
		return &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	return nil
}

func (trackUseCase *TrackUseCase) GetByArtist(id, amount int, isAuthorized bool) ([]models.Track, *models.CustomError) {
	tracks, err := trackUseCase.TrackRepository.GetByArtistID(id, amount, isAuthorized)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}

	return tracks, nil
}

func (trackUseCase *TrackUseCase) GetByAlbum(id int, isAuthorized bool) ([]models.Track, *models.CustomError) {
	tracks, err := trackUseCase.TrackRepository.GetByAlbumID(id, isAuthorized)
	if err != nil {
		fmt.Println(err)
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       constants.DatabaseNotResponding,
		}
	}

	return tracks, nil
}

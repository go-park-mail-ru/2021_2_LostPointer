package usecase

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/track"
	"2021_2_LostPointer/internal/utils/constants"
	"net/http"
)

type TrackUseCase struct {
	TrackRepository track.TrackRepository
}

func NewTrackUseCase(trackRepository track.TrackRepository) TrackUseCase {
	return TrackUseCase{TrackRepository: trackRepository}
}

func (trackUseCase TrackUseCase) GetHome(amount int, isAuthorized bool) ([]models.Track, *models.CustomError) {
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

func (trackUseCase TrackUseCase) IncrementListenCount(id int64) *models.CustomError {
	err := trackUseCase.TrackRepository.IncrementListenCount(id)
	if err != nil {
		return &models.CustomError{
			ErrorType: http.StatusInternalServerError,
			OriginalError: err,
		}
	}
	return nil
}

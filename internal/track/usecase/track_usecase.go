package usecase

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/track"
	"net/http"
)

const DatabaseNotResponding = "Database not responding"

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
			Message:       DatabaseNotResponding,
		}
	}

	return tracks, nil
}

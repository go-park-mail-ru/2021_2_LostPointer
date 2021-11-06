package track

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/track_usecase_mock.go -pkg mock . TrackUseCase:MockTrackUseCase
type TrackUseCase interface {
	GetHome(amount int, isAuthorized bool) ([]models.Track, *models.CustomError)
	IncrementListenCount(int64) *models.CustomError
	GetByArtist(id, amount int, isAuthorized bool) ([]models.Track, *models.CustomError)
	GetByAlbum(id int, isAuthorized bool) ([]models.Track, *models.CustomError)
}

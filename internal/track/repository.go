package track

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/track_repo_db_mock.go -pkg mock . TrackRepository:MockTrackRepository
type TrackRepository interface {
	GetRandom(amount int, isAuthorized bool) ([]models.Track, error)
	IncrementListenCount(int64) error
	GetByArtistID(artistID, amount int, isAuthorized bool) ([]models.Track, error)
	GetByAlbumID(albumID int, isAuthorized bool) ([]models.Track, error)
}

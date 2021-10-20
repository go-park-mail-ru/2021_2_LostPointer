package music

import (
	"2021_2_LostPointer/internal/models"
)

//go:generate moq -out ../mock/music_usecase_mock.go -internal mock . MusicUseCase:MockMusicUseCase
type MusicUseCase interface {
	GetMusicCollection(isAuthorized bool) (*models.MusicCollection, *models.CustomError)
	GetTracksForCollection(amount int, isAuthorized bool) ([]models.Track, error)
	GetAlbumsForCollection(amount int) ([]models.Album, error)
	GetArtistsForCollection(amount int) ([]models.Artist, error)
	GetPlaylistsForCollection(amount int) ([]models.Playlist, error)
}

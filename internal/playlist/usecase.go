package playlist

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/playlist_usecase_db_mock.go -pkg mock . PlaylistUseCase:MockPlaylistUseCase
type PlaylistUseCase interface {
	GetHome(amount int) ([]models.Playlist, *models.CustomError)
}


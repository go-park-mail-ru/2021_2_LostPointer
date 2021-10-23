package playlist

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/playlist_repo_db_mock.go -pkg mock . PlaylistRepository:MockPlayListRepository
type PlaylistRepository interface {
	GetRandom(amount int) ([]models.Playlist, error)
}

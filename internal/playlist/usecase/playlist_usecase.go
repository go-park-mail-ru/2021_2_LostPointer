package usecase

import (
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/playlist"
	"net/http"
)

const DatabaseNotResponding = "Database not responding"
const SiteID = 0

type PlaylistUseCase struct {
	PlaylistRepository playlist.PlaylistRepository
}

func NewPlaylistUseCase(playlistRepository playlist.PlaylistRepository) PlaylistUseCase {
	return PlaylistUseCase{PlaylistRepository: playlistRepository}
}

func (playlistUseCase PlaylistUseCase) GetHome(amount int) ([]models.Playlist, *models.CustomError) {
	playlists, err := playlistUseCase.PlaylistRepository.Get(amount, SiteID)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
			Message:       DatabaseNotResponding,
		}
	}

	return playlists, nil
}

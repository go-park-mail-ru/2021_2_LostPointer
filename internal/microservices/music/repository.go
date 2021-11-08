package music

import "2021_2_LostPointer/internal/models"

type MusicStorage interface {
	RandomTracks(int64, bool) ([]models.Track, error)
}
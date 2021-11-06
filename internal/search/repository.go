package search

import "2021_2_LostPointer/internal/models"

type MusicInfoStorage interface {
	TracksByFullWord(string) ([]models.Track, error)
	TracksByPartial(string) ([]models.Track, error)
	Artists(string) ([]models.ArtistShort, error)
}

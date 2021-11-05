package search

import "2021_2_LostPointer/internal/models"

type SearchRepository interface {
	SearchRelevantTracksByFullWord(string) ([]models.Track, error)
	SearchRelevantTracksByPartial(string) ([]models.Track, error)
	SearchRelevantArtists(string) ([]models.ArtistShort, error)
}

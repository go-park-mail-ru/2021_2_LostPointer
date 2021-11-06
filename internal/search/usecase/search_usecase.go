package usecase

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/search"
	"net/http"
)

type SearcherService struct {
	storage search.MusicInfoStorage
}

func NewSearchUsecase(storage search.MusicInfoStorage) SearcherService {
	return SearcherService{storage: storage}
}

func contains(tracks []models.Track, trackID int64) bool {
	for _, currentTrack := range tracks {
		if currentTrack.Id == trackID {
			return true
		}
	}
	return false
}

func (searcher *SearcherService) ByText(text string) (*models.SearchPayload, *models.CustomError) {
	relevantTracks, err := searcher.storage.TracksByFullWord(text)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
		}
	}
	amountRelevantTracksByFullWord := len(relevantTracks)

	var relevantTracksByPartial []models.Track
	if amountRelevantTracksByFullWord < constants.TracksSearchAmount {
		relevantTracksByPartial, err = searcher.storage.TracksByPartial(text)
		if err != nil {
			return nil, &models.CustomError{
				ErrorType:     http.StatusInternalServerError,
				OriginalError: err,
			}
		}
	}
	for _, track := range relevantTracksByPartial {
		if !contains(relevantTracks, track.Id) {
			relevantTracks = append(relevantTracks, track)
		}
	}

	relevantArtists, err := searcher.storage.Artists(text)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	searchData := &models.SearchPayload{
		Tracks:  relevantTracks,
		Artists: relevantArtists,
	}

	return searchData, nil
}

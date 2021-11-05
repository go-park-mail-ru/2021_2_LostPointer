package usecase

import (
	"2021_2_LostPointer/internal/constants"
	"2021_2_LostPointer/internal/models"
	"2021_2_LostPointer/internal/search"
	"net/http"
)

type SearchUsecase struct {
	searchDB search.SearchRepository
}

func NewSearchUsecase(searchDB search.SearchRepository) SearchUsecase {
	return SearchUsecase{searchDB: searchDB}
}

func contains(tracks []models.Track, trackID int64) bool {
	for _, currentTrack := range tracks {
		if currentTrack.Id == trackID {
			return true
		}
	}
	return false
}

func (searchU SearchUsecase) SearchByText(text string) (*models.SearchResponseData, *models.CustomError) {
	relevantTracks, err := searchU.searchDB.SearchRelevantTracksByFullWord(text)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
		}
	}
	amountRelevantTracksByFullWord := len(relevantTracks)

	var relevantTracksByPartial []models.Track
	if amountRelevantTracksByFullWord < constants.TracksSearchAmount {
		relevantTracksByPartial, err = searchU.searchDB.SearchRelevantTracksByPartial(text)
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

	relevantArtists, err := searchU.searchDB.SearchRelevantArtists(text)
	if err != nil {
		return nil, &models.CustomError{
			ErrorType:     http.StatusInternalServerError,
			OriginalError: err,
		}
	}

	searchData := &models.SearchResponseData{
		Tracks:  relevantTracks,
		Artists: relevantArtists,
	}

	return searchData, nil
}

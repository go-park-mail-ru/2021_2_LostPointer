package search

import "2021_2_LostPointer/internal/models"

type Searcher interface {
	ByText(string) (*models.SearchPayload, *models.CustomError)
}

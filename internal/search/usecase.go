package search

import "2021_2_LostPointer/internal/models"

type SearchUsecase interface {
	SearchByText(string) (*models.SearchResponseData, *models.CustomError)
}

package artist

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/artist_usecase_db_mock.go -pkg mock . ArtistUseCase:MockArtistUseCase
type ArtistUseCase interface {
	GetProfile(id int, isAuthorized bool) (*models.Artist, *models.CustomError)
}

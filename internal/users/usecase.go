package users

import (
	"2021_2_LostPointer/internal/models"
)

//go:generate moq -out ../mock/user_usecase_mock.go -pkg mock . UserUseCase:MockUserUseCase
type UserUseCase interface {
	Register(*models.User) (string, *models.CustomError)
	Login(*models.Auth) (string, *models.CustomError)
	Logout(string) error
	GetSettings(int) (*models.SettingsGet, *models.CustomError)
	UpdateSettings(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError
	GetAvatarFilename(int) (string, *models.CustomError)
}

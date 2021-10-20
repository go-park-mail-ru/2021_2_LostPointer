package users

import (
	"2021_2_LostPointer/internal/models"
)

//go:generate moq -out ../mock/user_usecase_mock.go -internal mock . UserUseCase:MockUserUseCase
type UserUseCase interface {
	Register(models.User) (string, *models.CustomError)
	Login(models.Auth) (string, *models.CustomError)
	IsAuthorized(string) (bool, int, *models.CustomError)
	Logout(string)
	GetSettings(int) (*models.SettingsGet, *models.CustomError)
	UpdateSettings(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError
}
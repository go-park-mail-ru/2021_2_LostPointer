package users

import (
	"2021_2_LostPointer/pkg/models"
)

//go:generate moq -out ../mock/user_usecase_mock.go -pkg mock . UserUseCaseIFace:MockUserUseCaseIFace
type UserUseCaseIFace interface {
	Register(models.User) (string, *models.CustomError)
	Login(models.Auth) (string, *models.CustomError)
	IsAuthorized(string) (bool, int)
	Logout(string)
	GetSettings(int) (*models.SettingsGet, *models.CustomError)
	UpdateSettings(int, *models.SettingsGet, *models.SettingsUpload) *models.CustomError
}

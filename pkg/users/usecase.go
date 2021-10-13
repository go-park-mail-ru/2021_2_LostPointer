package users

import "2021_2_LostPointer/pkg/models"

//go:generate moq -out ../mock/user_usecase_mock.go -pkg mock . UserUseCaseIFace:MockUserUseCaseIFace
type UserUseCaseIFace interface {
	Register(models.User) (string, string, error)
	Login(models.Auth) (string, error)
	Logout(string)
	IsAuthorized(string) (bool, error)
}

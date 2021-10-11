package users

import "2021_2_LostPointer/pkg/models"

type UserUseCaseIFace interface {
	Register(models.User) (string, string, error)
	Login(models.Auth) (string, error)
	GetSession(string) (bool, error)
	DeleteSession(string)
}

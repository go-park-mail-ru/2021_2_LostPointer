package users

import "2021_2_LostPointer/internal/models"

type UserRepository interface {
	CreateUser(models.User) (string, error)
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
}

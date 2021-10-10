package users

import "2021_2_LostPointer/internal/models"

//go:generate moq -out ../mock/user_repo_mock.go . UserRepository:MockUserRepository
type UserRepository interface {
	CreateUser(models.User) (uint64, error)
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	UserExits(models.Auth) (uint64, error)
}

package users

import "2021_2_LostPointer/pkg/models"

//go:generate moq -out ../mock/user_repo_mock.go -pkg mock . UserRepository:MockUserRepository
type UserRepositoryIFace interface {
	CreateUser(models.User, ...string) (uint64, error)
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	UserExits(models.Auth) (uint64, error)
}

type RedisStoreIFace interface {
	StoreSession(uint64) error
	GetSessionUserId(string) (int, error)
	DeleteSession(string)
}

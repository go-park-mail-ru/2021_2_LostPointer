package users

import "2021_2_LostPointer/pkg/models"

//go:generate moq -out ../mock/user_repo_db_mock.go -pkg mock . UserRepositoryIFace:MockUserRepositoryIFace
type UserRepositoryIFace interface {
	CreateUser(models.User, ...string) (uint64, error)
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	DoesUserExist(models.Auth) (uint64, error)
	GetSettings(int) (*models.Settings, error)
}

//go:generate moq -out ../mock/user_repo_redis_mock.go -pkg mock . RedisStoreIFace:MockRedisStoreIFace
type RedisStoreIFace interface {
	StoreSession(uint64, ...string) (string, error)
	GetSessionUserId(string) (int, error)
	DeleteSession(string)
}

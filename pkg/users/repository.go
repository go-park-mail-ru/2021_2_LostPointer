package users

import (
	"2021_2_LostPointer/pkg/models"
	"mime/multipart"
)

//go:generate moq -out ../mock/user_repo_db_mock.go -pkg mock . UserRepositoryIFace:MockUserRepositoryIFace
type UserRepositoryIFace interface {
	CreateUser(models.User, ...string) (uint64, error)
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	DoesUserExist(models.Auth) (uint64, error)
	GetSettings(int) (*models.SettingsGet, error)
	CheckPasswordByUserID(int, string) (bool, error)
	GetAvatarFilename(int) (string, error)
	UpdateEmail(int, string) error
	UpdateNickname(int, string) error
	UpdatePassword(int, string, ...string) error
	UpdateAvatar(int, string) error
}

//go:generate moq -out ../mock/user_repo_redis_mock.go -pkg mock . RedisStoreIFace:MockRedisStoreIFace
type RedisStoreIFace interface {
	StoreSession(uint64, ...string) (string, error)
	GetSessionUserId(string) (int, *models.CustomError)
	DeleteSession(string)
}

//go:generate moq -out ../mock/user_repo_filysystem_mock.go -pkg mock . FileSystemIFace:MockFileSystemIFace
type FileSystemIFace interface {
	CreateImage(*multipart.FileHeader) (string, error)
	DeleteImage(string) error
}

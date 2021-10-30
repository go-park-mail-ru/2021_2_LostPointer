package users

import (
	"2021_2_LostPointer/internal/models"
	"mime/multipart"
)

//go:generate moq -out ../mock/user_repo_db_mock.go -pkg mock . UserRepository:MockUserRepository
type UserRepository interface {
	CreateUser(*models.User, ...string) (int, error)
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	DoesUserExist(*models.Auth) (int, error)
	GetSettings(int) (*models.SettingsGet, error)
	CheckPasswordByUserID(int, string) (bool, error)
	GetAvatarFilename(int) (string, error)
	UpdateEmail(int, string) error
	UpdateNickname(int, string) error
	UpdatePassword(int, string, ...string) error
	UpdateAvatar(int, string) error
}

//go:generate moq -out ../mock/user_repo_filysystem_mock.go -pkg mock . FileSystem:MockFileSystem
type FileSystem interface {
	CreateImage(*multipart.FileHeader) (string, error)
	DeleteImage(string) error
}

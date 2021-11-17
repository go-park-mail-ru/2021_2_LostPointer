package profile

import (
	"2021_2_LostPointer/internal/models"
	"mime/multipart"
)

type UserSettingsStorage interface {
	GetSettings(int64) (*models.UserSettings, error)
	UpdateEmail(int64, string) error
	UpdateNickname(int64, string) error
	UpdatePassword(int64, string) error
	UpdateAvatar(int64, string) error
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	CheckPasswordByUserID(int64, string) (bool, error)
	CreateAvatar(*multipart.FileHeader) (string, error)
	DeleteAvatar(string) error
}

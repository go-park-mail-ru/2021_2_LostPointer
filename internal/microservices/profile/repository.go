package profile

import (
	"2021_2_LostPointer/internal/microservices/profile/proto"
)

//go:generate moq -out ./mock/profile_repo_mock.go -pkg mock . UserSettingsStorage:MockUserSettingsStorage
type UserSettingsStorage interface {
	GetSettings(int64) (*proto.UserSettings, error)
	UpdateEmail(int64, string) error
	UpdateNickname(int64, string) error
	UpdatePassword(int64, string) error
	UpdateAvatar(int64, string) error
	IsEmailUnique(string) (bool, error)
	IsNicknameUnique(string) (bool, error)
	CheckPasswordByUserID(int64, string) (bool, error)
}

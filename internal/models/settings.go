package models

type UserSettings struct {
	Email       string `json:"email,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	SmallAvatar string `json:"small_avatar,omitempty"`
	BigAvatar   string `json:"big_avatar,omitempty"`
}

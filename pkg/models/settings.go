package models

import "mime/multipart"

type SettingsGet struct {
	Email    	string `json:"email,omitempty"`
	Nickname 	string `json:"nickname,omitempty"`
	Avatar   	string `json:"avatar,omitempty"`
}

type SettingsUpload struct {
	Email    	   string `json:"email,omitempty"`
	Nickname 	   string `json:"nickname,omitempty"`
	Avatar   	   *multipart.FileHeader
	AvatarFileName string `json:"avatar,omitempty"`
	OldPassword    string `json:"old_password"`
	NewPassword    string `json:"new_password"`
}

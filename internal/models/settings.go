package models

import "mime/multipart"

type SettingsGet struct {
	Email       string `json:"email,omitempty"`
	Nickname    string `json:"nickname,omitempty"`
	SmallAvatar string `json:"small_avatar,omitempty"`
	BigAvatar   string `json:"big_avatar,omitempty"`
}

type SettingsUpload struct {
	Email    	   string `json:"email,omitempty" form:"email" query:"email"`
	Nickname 	   string `json:"nickname,omitempty" form:"nickname" query:"nickname"`
	Avatar   	   *multipart.FileHeader
	AvatarFileName string `json:"avatar,omitempty" form:"avatar" query:"avatar"`
	OldPassword    string `json:"old_password" form:"old_password" query:"old_password"`
	NewPassword    string `json:"new_password" form:"new_password" query:"new_password"`
}
